package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	Addr string
	env string
	db struct {
		dsn string
	}
}

type application struct {
	infoLog *log.Logger
	errorLog *log.Logger
	config config
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel();

	err = db.PingContext(ctx)
	if err != nil {
	return nil, err
	}
	return db, nil
}

func main() {

	var cfg config

	flag.StringVar(&cfg.Addr, "addr", ":4000", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()


	infoLog := log.New(os.Stdout, "INFO\t", log.Default().Flags())
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLog: infoLog,
		errorLog: errorLog,
		config: cfg,
	}

	db, err := openDB(cfg)
	if err != nil {
		errorLog.Fatal()
	}

	defer db.Close()

	infoLog.Printf("database connection pool established")

	srv := &http.Server{
		ErrorLog: app.errorLog,
		Addr: cfg.Addr,
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	infoLog.Printf("Starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)


}