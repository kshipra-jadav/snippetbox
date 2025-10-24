package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kshipra-jadav/snippetbox/internal/models"
	"github.com/kshipra-jadav/snippetbox/internal/validator"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}
type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *App) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

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

	app.sessionManager.GetString(r.Context(), "flash")

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.html", data)
}

func (app *App) snippetCreateGet(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("isAuthenticated --- >", "auth status", app.isAuthenticated(r.Context()))
	data := app.newTemplateData(r)
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
		app.serverError(w, r, err)
		return
	}

	myForm.CheckField(validator.NotBlank(myForm.Title), "title", "Title field cannot be blank.")
	myForm.CheckField(validator.MaxChars(myForm.Title, 100), "title", "Title field has to be less than 100 chars.")

	myForm.CheckField(validator.NotBlank(myForm.Content), "content", "Content field cannot be blank")

	myForm.CheckField(validator.PermittedValue(myForm.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365.")

	if !myForm.Valid() {
		data := app.newTemplateData(r)
		data.Form = myForm
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

func (app *App) userSignupGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = userSignupForm{}

	app.render(w, r, http.StatusOK, "signup.html", data)

}

func (app *App) userSignupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var signupForm userSignupForm
	err = app.formDecoder.Decode(&signupForm, r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	signupForm.Validator.CheckField(validator.NotBlank(signupForm.Name), "name", "name cannot be empty")
	signupForm.Validator.CheckField(validator.NotBlank(signupForm.Email), "email", "email cannot be empty")
	signupForm.Validator.CheckField(validator.NotBlank(signupForm.Password), "password", "password cannot be empty")

	signupForm.Validator.CheckField(validator.ValidEmail(signupForm.Email), "email", "email must be in correct format")
	signupForm.Validator.CheckField(validator.MinChars(signupForm.Password, 8), "password", "password should be greater than 8 characters")

	if !signupForm.Valid() {
		data := app.newTemplateData(r)
		data.Form = signupForm
		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(signupForm.Name, signupForm.Email, signupForm.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			signupForm.AddFieldError("email", "Email already in use")
			data := app.newTemplateData(r)
			data.Form = signupForm
			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
			return
		} else {
			app.serverError(w, r, err)
			return
		}
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please login!")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

func (app *App) userLoginGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = userSignupForm{}

	app.render(w, r, http.StatusOK, "login.html", data)
}

func (app *App) userLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var loginForm userLoginForm
	err = app.formDecoder.Decode(&loginForm, r.PostForm)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	loginForm.CheckField(validator.NotBlank(loginForm.Email), "email", "email cannot be blank")
	loginForm.CheckField(validator.NotBlank(loginForm.Password), "password", "password cannot be blank")
	loginForm.CheckField(validator.ValidEmail(loginForm.Email), "email", "email is not valid")

	if !loginForm.Valid() {
		data := app.newTemplateData(r)
		data.Form = loginForm
		app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	usrId, err := app.users.Authenticate(loginForm.Email, loginForm.Password)
	if err != nil {
		if errors.Is(err, models.ErrNoRecords) {
			loginForm.AddNonFieldError("No user found with this email. Please try again.")
			data := app.newTemplateData(r)
			data.Form = loginForm
			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
			return
		} else if errors.Is(err, models.ErrInvalidCredentials) {
			loginForm.AddNonFieldError("Invalid username or password. Please try again")
			data := app.newTemplateData(r)
			data.Form = loginForm
			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
			return
		} else {
			app.serverError(w, r, err)
			return
		}
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserId", usrId)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *App) userLogout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserId")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
