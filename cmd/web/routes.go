package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *App) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("../../ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{snippetID}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreateGet))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	return alice.New(app.recoverPanic, app.logRequest).Then(commonHeaders(mux))
}
