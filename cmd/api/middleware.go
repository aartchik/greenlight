package main

import (
	"fmt"
	"net/http"
	"golang.org/x/time/rate"
	"sync"
	"net"
	"time"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		app.logger.PrintInfo(fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI()), nil)

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {

	type client struct {
		limiter *rate.Limiter
		LastSeen time.Time
	}
	var (
		mu sync.Mutex
		clients = make(map[string]*client)
		)

		go func() {
			for {
				time.Sleep(time.Minute) 
				mu.Lock()
				for ip, client := range clients {
					if time.Since(client.LastSeen) > 3 * time.Minute {
						delete(clients, ip)
					}
				}
				mu.Unlock()
			}
		} ()
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return 
		}
		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
			}
		} 
		clients[ip].LastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return 
		}

		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Closed")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
				return
			}
		} ()

		next.ServeHTTP(w, r)
	})
}