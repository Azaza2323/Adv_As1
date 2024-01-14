package main

import (
	mysql2 "asik1/pkg/mysql"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

type application struct {
	debug    bool
	errorLog *log.Logger
	infoLog  *log.Logger
	news     *mysql2.NewsModel
}

func main() {
	addr := flag.String("addr", ":3000", "HTTP network address")
	debug := flag.Bool("debug", false, "Enable debug mode")
	dsn := flag.String("dsn", "root:WheelyStrongPwd@/articles?parseTime=true", "news")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	var model *mysql2.NewsModel
	model = &mysql2.NewsModel{DB: db}

	app := &application{
		debug:    *debug,
		errorLog: errorLog,
		infoLog:  infoLog,
		news:     model,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
