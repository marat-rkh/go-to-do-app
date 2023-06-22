package main

import (
	"fmt"
	"log"
	"net/http"

	"server/middleware"
	"server/router"
)

func main() {
	middleware.InitTaskRegistry()
	r := router.Router()
	// fs := http.FileServer(http.Dir("build"))
	// http.Handle("/", fs)
	fmt.Println("Starting server on the port 8080...")

	log.Fatal(http.ListenAndServe(":8080", r))
}
