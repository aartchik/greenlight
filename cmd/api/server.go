package main
import (
	"net/http"
	"time"
)

func (app *application) serve() error {

	srv := &http.Server{
		Addr: 			app.config.port,
		Handler:		app.routes(),
		IdleTimeout:	time.Minute,
		ReadTimeout: 	10 * time.Second,
		WriteTimeout: 	30 * time.Second,
	}
	app.logger.PrintInfo("starting server", map[string]string{ "addr": srv.Addr, "env": app.config.env,})

	return srv.ListenAndServe()
}