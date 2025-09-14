package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/kshipra-jadav/snippetbox/internal/models"
	"github.com/kshipra-jadav/snippetbox/internal/validator"
)

type App struct {
	logger         *slog.Logger
	snippets       *models.SnippetsModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
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

	flashMsg := app.sessionManager.GetString(r.Context(), "flash")

	data := templateData{
		Snippet: snippet,
		Flash:   flashMsg,
	}

	app.render(w, r, http.StatusOK, "view.html", data)
}

func (app *App) snippetCreateGet(w http.ResponseWriter, r *http.Request) {
	data := templateData{}
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, r, http.StatusOK, "create.html", data)
}

func (app *App) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var myForm snippetCreateForm

	err = app.formDecoder.Decode(&myForm, r.PostForm)
	if err != nil {
		fmt.Println("LMAOOOOO")
		app.serverError(w, r, err)
		return
	}

	myForm.CheckField(validator.NotBlank(myForm.Title), "title", "Title field cannot be blank.")
	myForm.CheckField(validator.MaxChars(myForm.Title, 100), "title", "Title field has to be less than 100 chars.")

	myForm.CheckField(validator.NotBlank(myForm.Content), "content", "Content field cannot be blank")

	myForm.CheckField(validator.PermittedValue(myForm.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365.")

	fmt.Println("fieldErrors", myForm.FieldErrors)

	if !myForm.Valid() {
		data := templateData{Form: myForm}
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	lastID, err := app.snippets.Insert(myForm.Title, myForm.Content, myForm.Expires)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Snippet with ID: %v, written successfully.", lastID))

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%v", lastID), http.StatusSeeOther)
}
