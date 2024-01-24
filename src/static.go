package src

import (
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/alice"
)

type Static struct {
	route       string
	filepath    string
	contentType string
}

// var statics = []Static{
// 	{"/static/css/tw-output.css", "src/static/css/tw-output.css", "text/css"},
// 	{"/static/js/alert.js", "src/static/js/alert.js", "text/javascript"},
// 	{"/favicon.ico", "src/static/favicon.ico", ""},
// 	{"/white-bear.ico", "src/static/white-bear.ico", ""},
// }

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

		asdf := path[10:]
		log.Printf("asdf: %v, path: %v\n", asdf, path)

		if jsFile, err := filepath.Match("*.js", filepath.Base(path)); jsFile {
			files = append(files, Static{route: asdf, filepath: path, contentType: "text/javascript"})
			return err
		}
		if cssFile, err := filepath.Match("*.css", filepath.Base(path)); cssFile {
			files = append(files, Static{route: asdf, filepath: path, contentType: "text/css"})
			return err
		}
		if imgFile, err := filepath.Match("*.{png, ico, jpg, webm}", filepath.Base(path)); imgFile {
			files = append(files, Static{route: asdf, filepath: path, contentType: ""})
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
	}

}
