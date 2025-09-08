package main

import "net/http"

func (app *App) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error(), "method:", r.Method, "URI", r.URL.RequestURI())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *App) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
