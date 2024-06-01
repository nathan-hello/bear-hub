package src

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/justinas/alice"
	"github.com/nathan-hello/htmx-template/src/routes"
)

func SiteRouter() {

	type Site struct {
		route       string
		hfunc       http.HandlerFunc
		middlewares alice.Chain
	}

	sites := []Site{
                {route: "/",
			hfunc: routes.Root,
			middlewares: alice.New(
				RejectSubroute("/"),
				Logging,
				AllowMethods("GET"),
				InjectClaimsOnValidToken,
			)},
		{route: "/todo",
			hfunc: routes.Todo,
			middlewares: alice.New(
				Logging,
				AllowMethods("GET", "DELETE", "POST"),
				InjectClaimsOnValidToken,
			)},
		{route: "/auth",
			hfunc: routes.Auth,
			middlewares: alice.New(
				Logging,
				AllowMethods("GET"),
				InjectClaimsOnValidToken,
			)},
		{route: "/auth/signup",
			hfunc: routes.SignUp,
			middlewares: alice.New(
				Logging,
				AllowMethods("GET", "POST"),
				InjectClaimsOnValidToken,
			)},
		{route: "/auth/signout",
			hfunc: routes.SignOut,
			middlewares: alice.New(
				Logging,
				AllowMethods("GET", "POST"),
				InjectClaimsOnValidToken,
			)},
		{route: "/auth/signin",
			hfunc: routes.SignIn,
			middlewares: alice.New(
				Logging,
				AllowMethods("GET", "POST"),
				InjectClaimsOnValidToken,
			)},
		{route: "/profile/",
			hfunc: routes.UserProfile,
			middlewares: alice.New(
				Logging,
				AllowMethods("GET"),
				InjectClaimsOnValidToken,
			)},
		{route: "/chat",
			hfunc: routes.Chat,
			middlewares: alice.New(
				Logging,
				AllowMethods("GET", "POST", "DELETE", "PUT"),
				InjectClaimsOnValidToken,
			)},
		{route: "/ws/v1/chat/html",
			hfunc: routes.ChatSocket,
			middlewares: alice.New(
				Logging,
			)},
		{route: "/api/v1/chat/message",
			hfunc: routes.ApiChat,
			middlewares: alice.New(
				AllowMethods("POST"),
			)},
	}

	for _, v := range sites {
		http.Handle(v.route, v.middlewares.ThenFunc(v.hfunc))
	}

	http.ListenAndServe(":3001", nil)
}

type Static struct {
	route       string
	filepath    string
	contentType string
}

func LoadStaticFiles() ([]Static, error) {

	files := []Static{}
	publicDir := "public/"
	images := []string{".ico", ".png", ".jpg", ".webm"}
	plain := []string{".js", ".css"}
	allowed := append(images, plain...)

	err := filepath.Walk("public", func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(info.Name())
		match := slices.Contains(allowed, ext)
		if !match {
			fmt.Printf("file %v not in allow list", path)
			return nil
		}

		staticRoute := strings.TrimPrefix(path, publicDir)
		staticRoute = "/" + staticRoute

		jsFile, err := filepath.Match("*.js", filepath.Base(path))
		if jsFile {
			files = append(files, Static{route: staticRoute, filepath: path, contentType: "text/javascript"})
			return err
		}
		cssFile, err := filepath.Match("*.css", filepath.Base(path))
		if cssFile {
			files = append(files, Static{route: staticRoute, filepath: path, contentType: "text/css"})
			return err
		}
		files = append(files, Static{route: staticRoute, filepath: path, contentType: ""})
		return err
	})

	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no static files: %#v", files)
	}

	return files, nil
}

func StaticRouter(files []Static) error {
	for _, v := range files {
		middles := alice.New(CreateHeader("Content-Type", v.contentType))
		file := v.filepath // closure shanigans
		http.Handle(v.route, middles.ThenFunc(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, file) }))
		log.Printf("Creating route: %v, for file: %v, with Content-Type %v\n", v.route, file, v.contentType)
	}
	return nil

}
