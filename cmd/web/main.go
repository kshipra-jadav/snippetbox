package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// snippet view, snippet create, snippet create post, home
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{snippetID}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreateGet)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	fmt.Println("Listening on localhost:4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatalf("Error : %v", err)
}
