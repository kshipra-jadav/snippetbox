package main

import (
	"github.com/kshipra-jadav/snippetbox/ui"
	"net/http"

	"github.com/justinas/alice"
)

func (app *App) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServerFS(ui.Files)
	mux.Handle("GET /static/", fileServer)

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{snippetID}", dynamic.ThenFunc(app.snippetView))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignupGet))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))

	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLoginGet))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreateGet))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogout))

	return alice.New(app.recoverPanic, app.logRequest).Then(commonHeaders(mux))
}
