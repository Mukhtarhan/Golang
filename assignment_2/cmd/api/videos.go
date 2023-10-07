package main

import (
	"fmt"
	"net/http"
	"time"

	"assignment_2.alexedwards.net/internal/data"
)

// Add a createVideoHandler for the "POST /v1/Videos" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createVideoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "add a new video")
}

// Add a showVideoHandler for the "GET /v1/Videos/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showVideoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		// Use the new notFoundResponse() helper.
		app.notFoundResponse(w, r)
		return
	}
	video := data.Video{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"Video": video}, nil)
	if err != nil {
		// Use the new serverErrorResponse() helper.
		app.serverErrorResponse(w, r, err)
	}
}
