package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there. This is the home!")
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	snippetID, err := strconv.Atoi(r.PathValue("snippetID"))
	if err != nil || snippetID <= 0 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "You're viewing snippet number: %v", snippetID)
}

func snippetCreateGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is form for snippet creation")
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Snippet created successfully!")
}
