package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *App) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("../../ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{snippetID}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreateGet)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	return alice.New(app.recoverPanic, app.logRequest).Then(commonHeaders(mux))
}
