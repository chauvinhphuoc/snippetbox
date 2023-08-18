package main

import (
	"database/sql"
	"github.com/chauvinhphuoc/snippetbox/internal/db/sqlc"
	"github.com/go-playground/form/v4"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var (
	infoLog        = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog       = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	dataSourceName = "postgres://postgres:12345@localhost:5432/snippetbox?sslmode=disable"
)

// Everything inside an application is called an application's dependency,
// it sticks to the application for doing tasks.
type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	db       *sql.DB
	*sqlc.Queries
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder // A Decoder instance is used to map HTML field values into struct fields.
}

func main() {
	db, err := openDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Print("Connected to database")

	q := sqlc.New(db)

	templateCache, err := initialTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		db:            db,
		Queries:       q,
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	server := &http.Server{
		Addr:    "127.0.0.1:4000",
		Handler: app.routes(),
	}

	infoLog.Print("Starting server on http://localhost:4000")
	err = server.ListenAndServe()
	errorLog.Print(err)
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	if pingErr := db.Ping(); pingErr != nil {
		return nil, err
	}

	return db, nil
}
