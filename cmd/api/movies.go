package main

import (
	"fmt"
	"greenlight.aartchik.net/internal/data"
	"net/http"
	"time"
	"greenlight.aartchik.net/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	var input struct{
		Title     string 		`json:"title"`
		Year      int32			`json:"year"`
		Runtime   data.Runtime	`json:"runtime"`
		Genres    []string 		`json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	movies := &data.Movie{
		Title: input.Title,
		Year: input.Year,
		Runtime: input.Runtime,
		Genres: input.Genres,
	}

	data.ValidateMovie(v, movies)
	if data.ValidateMovie(v, movies); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}



	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	movie := data.Movie{
		Id:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r ,err)
	}
	
}
