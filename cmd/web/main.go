package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"snippetbox.essa/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	snippets *models.SnippetModel
}

func main() {

	//addr is a command line argument
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:essa28=1@/snippetbox?parseTime=true", "SQL")
	flag.Parse()

	//leled loggig
	//infoLog is used to reperesent information about the programm eg:log.printf("starting in :40000")
	infoLog := log.New(os.Stdout, " INFO ", log.Ldate|log.Ltime)
	//errorLog is used to show errors log.Lshortfile shows the file name and the line number
	errorLog := log.New(os.Stderr, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &models.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	app.infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
