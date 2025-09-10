package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/kshipra-jadav/snippetbox/internal/models"
)

type App struct {
	logger        *slog.Logger
	snippets      *models.SnippetsModel
	templateCache map[string]*template.Template
}

func (app *App) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		Snippets: snippets,
	}

	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *App) snippetView(w http.ResponseWriter, r *http.Request) {
	snippetID, err := strconv.Atoi(r.PathValue("snippetID"))
	if err != nil || snippetID <= 0 {
		app.logger.Error(err.Error())
		http.NotFound(w, r)
		return
	}
	snippet, err := app.snippets.Get(snippetID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		Snippet: snippet,
	}

	app.render(w, r, http.StatusOK, "view.html", data)
}

func (app *App) snippetCreateGet(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"../../ui/html/pages/base.html",
		"../../ui/html/pages/create.html",
		"../../ui/html/pages/footer.html",
	}

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

func (app *App) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	lastID, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Snippet with id: %v - Created Successfully!", lastID)
}
