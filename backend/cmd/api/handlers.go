package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type envelope map[string]interface{}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {

	type credentials struct {
		UserName string `json:"email"`
		Password string `json:"password"`
	}

	var creds credentials
	var payload jsonResponse

	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorMessage(w, err)
		return
	}
	user, err := app.models.User.GetByEmail(creds.UserName)
	if err != nil {
		app.errorMessage(w, errors.New("invalid username/password 1"))
		return
	}
	validPassword, err := user.PasswordMatches(creds.Password)
	if err != nil || !validPassword {
		app.errorMessage(w, errors.New("invalid username/password 2"))
		return
	}
	token, err := app.models.Token.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		app.errorMessage(w, err)
		return
	}
	err = app.models.Token.Insert(*token, *user)
	if err != nil {
		app.errorMessage(w, err)
		return
	}
	payload = jsonResponse{
		Error:   false,
		Message: "logged in",
		Data:    envelope{"token": token, "user": user},
	}
	if err = app.writeJSON(w, 200, payload); err != nil {
		log.Fatal("Fehler...")
	}
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Token string `json:"token"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorMessage(w, errors.New("invalid json"))
		return
	}
	err = app.models.Token.DeleteByToken(requestPayload.Token)
	if err != nil {
		app.errorMessage(w, errors.New("error in database"), 500)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "logged out",
	}

	err = app.writeJSON(w, 200, payload)
	if err != nil {
		app.errorMessage(w, errors.New("invalid json"))
	}
}
