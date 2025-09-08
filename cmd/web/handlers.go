package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

type App struct {
	logger *slog.Logger
}

func (app *App) home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"../../ui/html/pages/base.html",
		"../../ui/html/pages/home.tmpl",
		"../../ui/html/pages/footer.html",
	}
	app.logger.Info("Received request.", "method", r.Method, "URI", r.URL.RequestURI())
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *App) snippetView(w http.ResponseWriter, r *http.Request) {
	snippetID, err := strconv.Atoi(r.PathValue("snippetID"))
	if err != nil || snippetID <= 0 {
		app.logger.Error(err.Error())
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "You're viewing snippet number: %v", snippetID)
}

func (app *App) snippetCreateGet(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"../../ui/html/pages/base.html",
		"../../ui/html/pages/create.html",
		"../../ui/html/pages/footer.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		app.logger.Error(err.Error(), "method:", r.Method, "URI", r.URL.RequestURI())
		http.Error(w, "Internal Server Error Parsing The Template File", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		app.logger.Error(err.Error(), "method:", r.Method, "URI", r.URL.RequestURI())
		http.Error(w, "Internal Server Error Executing the Template File", http.StatusInternalServerError)
		return
	}
}

func (app *App) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Snippet created successfully!")
}
