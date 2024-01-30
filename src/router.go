package src

import (
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/alice"
	"github.com/nathan-hello/htmx-template/src/routes"
)

type Static struct {
	route       string
	filepath    string
	contentType string
}

func staticGet(filepath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath)
	}
}

func HandlePublic() {

	files := []Static{}

	filepath.Walk("src/public", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		staticRoute := path[10:]

		if jsFile, err := filepath.Match("*.js", filepath.Base(path)); jsFile {
			files = append(files, Static{route: staticRoute, filepath: path, contentType: "text/javascript"})
			return err
		}
		if cssFile, err := filepath.Match("*.css", filepath.Base(path)); cssFile {
			files = append(files, Static{route: staticRoute, filepath: path, contentType: "text/css"})
			return err
		}
		imgExts := []string{".ico", ".png", ".jpg", ".webm"}

		for _, v := range imgExts {
			match, err := filepath.Match("*"+v, filepath.Base(path))
			if err != nil {
				return err
			}
			if !match {
				continue
			}
			files = append(files, Static{route: staticRoute, filepath: path, contentType: ""})
			return err
		}

		return err
	})

	defaultMiddleware := alice.New(Logging, AllowMethods("GET"))
	for _, v := range files {
		if v.contentType != "" {
			defaultMiddleware.Append(CreateHeader("Content-Type", v.contentType))
		}
		http.Handle(v.route, defaultMiddleware.ThenFunc(staticGet(v.filepath)))
		log.Printf("Creating route: %v, for file: %v, with Content-Type %v\n", v.route, v.filepath, v.contentType)
	}

}

func HandleSites() {

	type Site struct {
		route       string
		hfunc       http.HandlerFunc
		middlewares alice.Chain
	}

	sites := []Site{
		{route: "/",
			hfunc: routes.Root,
			middlewares: alice.New(
				Logging,
				InjectClaimsOnValidToken,
				AllowMethods("GET"),
				RejectSubroute("/"),
			)},
		{route: "/todo",
			hfunc: routes.Todo,
			middlewares: alice.New(
				Logging,
				InjectClaimsOnValidToken,
				AllowMethods("GET", "DELETE", "POST"),
			)},
		{route: "/signup",
			hfunc: routes.SignUp,
			middlewares: alice.New(
				Logging,
				InjectClaimsOnValidToken,
				AllowMethods("GET", "POST"),
			)},
		{route: "/signin",
			hfunc: routes.SignIn,
			middlewares: alice.New(
				Logging,
				InjectClaimsOnValidToken,
				AllowMethods("GET", "POST"),
			)},
		{route: "/profile",
			hfunc: routes.UserProfile,
			middlewares: alice.New(
				Logging,
				InjectClaimsOnValidToken,
				AllowMethods("GET"),
			)},
		{route: "/c",
			hfunc: routes.MicroComponents,
			middlewares: alice.New(
				Logging,
				AllowMethods("GET"),
			)},
		{route: "/chat",
			hfunc: routes.ChatRoot,
			middlewares: alice.New(
				Logging,
				AllowMethods("GET", "POST", "DELETE", "PUT"),
			)},
		{route: "/chat/",
			hfunc: routes.ChatSubRouter,
			middlewares: alice.New(
				Logging,
				AllowMethods("GET", "POST", "DELETE", "PUT"),
			)},
	}

	for _, v := range sites {
		http.Handle(v.route, v.middlewares.ThenFunc(v.hfunc))
	}

	http.ListenAndServe(":3000", nil)
}
