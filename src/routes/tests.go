package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func say(w http.ResponseWriter, s string) {
	w.Write([]byte(fmt.Sprintf("%#v\n", s)))
}

func sampleProtectedRoute(w http.ResponseWriter, r *http.Request) {

	say(w, "hello from protected route")
	err := utils.VerifyJWT(r.Header["Token"][0])

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	say(w, "a big secret! i really like this girl named natalie!")
}

func signUp(w http.ResponseWriter, r *http.Request) {
	access, err := utils.CreateAccess(uuid.New())

	if err != nil {
		say(w, "error in signUproute")
		panic(err)
	}

	w.Write([]byte(access))
}

func TestingEntryPoint(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte(fmt.Sprintf("%#v\n", r.URL)))
	say(w, fmt.Sprintf("%v", time.Now().Unix()))
	if r.URL.Path == "/testing/signup" {
		signUp(w, r)
	}

	if r.URL.Path == "/testing/protected" {
		say(w, "meow")
		sampleProtectedRoute(w, r)
	}

}
