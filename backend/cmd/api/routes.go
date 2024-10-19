package main

import (
	"gowebapp/internal/data"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*", "https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	mux.Route("/users", func(r chi.Router) {
		r.Post("/login", app.Login)
		r.Post("/logout", app.Logout)

		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			var users data.User
			all, err := users.GetAll()
			if err != nil {
				log.Fatal(err)
			}
			payload := jsonResponse{
				Error:   false,
				Message: "success",
				Data:    envelope{"users": all},
			}

			err = app.writeJSON(w, 200, payload)
			if err != nil {
				log.Fatal("Fehler in DB")
			}
		})
		r.Get("/add", func(w http.ResponseWriter, r *http.Request) {
			var u = data.User{
				Email:     "janina@web.ch",
				FirstName: "Janina",
				LastName:  "Schmid",
				Password:  "freizeit",
			}
			app.infoLog.Println("Adding user...")

			id, err := app.models.User.Insert(u)
			if err != nil {
				app.errorLog.Println(err)
				app.errorMessage(w, err, http.StatusForbidden)
				return
			}
			app.infoLog.Println("Go back id of", id)
		})
		r.Get("/test-generate-token", func(w http.ResponseWriter, r *http.Request) {
			token, err := app.models.User.Token.GenerateToken(1, 60*time.Minute)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			token.Email = "admin@web.ch"
			token.CreatedAt = time.Now()
			token.UpatedAt = time.Now()

			payload := jsonResponse{
				Error:   false,
				Message: "success",
				Data:    token,
			}
			app.writeJSON(w, 200, payload)
		})
		r.Get("/test-save-token", func(w http.ResponseWriter, r *http.Request) {
			token, err := app.models.User.Token.GenerateToken(1, 60*time.Minute)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			user, err := app.models.User.GetById(1)
			if err != nil {
				app.errorLog.Println(err, "erst")
				return
			}
			token.UserID = user.ID
			token.Email = user.Email
			token.CreatedAt = time.Now()
			token.UpatedAt = time.Now()

			err = token.Insert(*token, *user)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			payload := jsonResponse{
				Error:   false,
				Message: "success",
				Data:    token,
			}
			app.writeJSON(w, 200, payload)
		})
		r.Get("/validating", func(w http.ResponseWriter, r *http.Request) {
			tokenToValidate := r.URL.Query().Get(("token"))
			valid, err := app.models.Token.ValidToken(tokenToValidate)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			payload := jsonResponse{
				Error:   false,
				Message: "success",
				Data:    valid,
			}
			app.writeJSON(w, 200, payload)
		})
	})
	mux.Route("/admin", func(r chi.Router) {
		r.Use(app.AuthTokenMiddleware)

		r.Post("/foo", func(w http.ResponseWriter, r *http.Request) {
			payload := jsonResponse{
				Error:   false,
				Message: "bar",
			}
			app.writeJSON(w, 200, payload)
		})
	})

	return mux
}
