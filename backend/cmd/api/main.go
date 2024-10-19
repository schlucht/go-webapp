package main

import (
	"fmt"
	"gowebapp/internal/data"
	"gowebapp/internal/driver"
	"log"
	"net/http"
	"os"
)

type config struct {
	port int
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	models   data.Models
	environment string
}

func main() {
	var cfg config
	cfg.port = 8081

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime)

	// dsn := "host=localhost port=5432 user=postgres password=password dbname=vueapi sslmode=disable timezone=UTC connect_timeout=5"
	dsn := os.Getenv("DSN")
	environment := os.Getenv("ENV")
	db, err := driver.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect db")
	}
	defer db.SQL.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		models:   data.New(db.SQL),
		environment: environment,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) serve() error {

	app.infoLog.Println("API listeing on port", app.config.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}
	return srv.ListenAndServe()
}
