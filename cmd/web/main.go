package main

import (
	"asik1/pkg/mysql"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	debug       bool
	errorLog    *log.Logger
	infoLog     *log.Logger
	session     *sessions.Session
	news        *mysql.NewsModel
	departments *mysql.DepoModel
	users       *mysql.UserModel
}

func main() {
	addr := flag.String("addr", ":1111", "HTTP network address")
	debug := flag.Bool("debug", false, "Enable debug mode")
	dsn := flag.String("dsn", "root:123123@/data_for_news?parseTime=true", "news.html")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	app := &application{
		debug:       *debug,
		errorLog:    errorLog,
		infoLog:     infoLog,
		session:     session,
		news:        &mysql.NewsModel{DB: db},
		departments: &mysql.DepoModel{DB: db},
		users:       &mysql.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
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
