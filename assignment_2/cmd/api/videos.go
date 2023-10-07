package main

import (
	"fmt"
	"net/http"
	"time"

	"assignment_2.alexedwards.net/internal/data"
	"assignment_2.alexedwards.net/internal/validator"
)

// Add a createVideoHandler for the "POST /v1/Videos" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createVideoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the input struct to a new video struct.
	video := &data.Video{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}
	// Initialize a new Validator.
	v := validator.New()
	// Call the Validatevideo() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateVideo(v, video); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)

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
