package main

import (
	"aitu.edu.kz/aitu-news/conf"
	"aitu.edu.kz/aitu-news/pkg/models/psql"
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
)

type application struct {
	config       *conf.Config
	errorLog     *log.Logger
	infoLog      *log.Logger
	articles     *psql.ArticleModel
	users        *psql.UserModel
	sessionStore *sessions.CookieStore
}

var app *application

func main() {
	//custom loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//reading config file
	cnf := conf.GetConfig()

	//database connection
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable ",
		cnf.DB.Host, cnf.DB.Port, cnf.DB.User, cnf.DB.Password, cnf.DB.DbName)
	db, err := openDB(psqlInfo)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	//global application variables
	app = &application{
		config:       conf.GetConfig(),
		errorLog:     errorLog,
		infoLog:      infoLog,
		articles:     &psql.ArticleModel{DB: db},
		users:        &psql.UserModel{DB: db},
		sessionStore: sessions.NewCookieStore([]byte(cnf.SERVER.SessionSecret)),
	}

	//creating server
	mux := http.NewServeMux()
	for route, handler := range *getRoutes() {
		mux.HandleFunc(route, handler)
	}
	fileServer := http.FileServer(http.Dir(cnf.SERVER.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	srv := &http.Server{
		Addr:     ":" + cnf.SERVER.Port,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	infoLog.Println("Starting server on :" + cnf.SERVER.Port)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}
