package main

import (
	"log"

	"github.com/nathan-hello/htmx-template/src"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func main() {
	err := utils.InitEnv()
	if err != nil {
		log.Fatal(err)
	}
	src.PublicRouter()
	src.SiteRouter()
}
