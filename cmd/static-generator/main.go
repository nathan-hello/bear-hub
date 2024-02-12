package main

import (
	"fmt"
	"net/http"
)

func main() {
	asdf, err := http.Get("http://localhost:3000")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", asdf)
}
