package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	_ "github.com/lib/pq"
	"greenlight.aartchik.net/internal/data"
	"greenlight.aartchik.net/internal/jsonlog"
)

const version = "1.0.0"

type config struct {
	port string
	env string
	db struct {
		dsn string
	}
	limiter struct {
		rps float64
		burst int
		enabled bool
	}
}

type application struct {
	logger *jsonlog.Logger
	config config
	models data.Models
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

	flag.StringVar(&cfg.port, "port", ":4000", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("db_dsn"), "PostgreSQL DSN")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	app := &application{
		logger: logger,
		config: cfg,
		models: data.NewModels(db),
	}



	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)


	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}