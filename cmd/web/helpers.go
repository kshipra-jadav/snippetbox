package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *App) serverError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Println(string(debug.Stack()))
	app.logger.Error(err.Error(), "method:", r.Method, "URI", r.URL.RequestURI())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *App) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *App) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	tmpl, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %v doesn't exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)

}
