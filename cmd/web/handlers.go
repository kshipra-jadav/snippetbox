package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"../../ui/html/pages/base.html",
		"../../ui/html/pages/home.tmpl",
		"../../ui/html/pages/footer.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error Parsing The Template File", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error Executing the Template File", http.StatusInternalServerError)
		return
	}
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
	files := []string{
		"../../ui/html/pages/base.html",
		"../../ui/html/pages/create.html",
		"../../ui/html/pages/footer.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error Parsing The Template File", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error Executing the Template File", http.StatusInternalServerError)
		return
	}
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Snippet created successfully!")
}
