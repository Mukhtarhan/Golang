package main

import (
	"errors"
	"fmt"
	"net/http"

	"assignment_2.alexedwards.net/internal/data"
	"assignment_2.alexedwards.net/internal/validator"
	// "time"
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
	err = app.models.Videos.Insert(video)

	if data.ValidateVideo(v, video); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/videos/%d", video.ID))
	// Write a JSON response with a 201 Created status code, the video data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"video": video}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	fmt.Fprintf(w, "%+v\n", input)

}

// Add a showVideoHandler for the "GET /v1/Videos/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showVideoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Call the Get() method to fetch the data for a specific video. We also need to
	// use the errors.Is() function to check if it returns a data.ErrRecordNotFound
	// error, in which case we send a 404 Not Found response to the client.
	video, err := app.models.Videos.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"video": video}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateVideoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the video ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the existing video record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	video, err := app.models.Videos.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Declare an input struct to hold the expected data from the client.
	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}
	// Decode the JSON as normal.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// If the input.Title value is nil then we know that no corresponding "title" key/
	// value pair was provided in the JSON request body. So we move on and leave the
	// video record unchanged. Otherwise, we update the video record with the new title
	// value. Importantly, because input.Title is a now a pointer to a string, we need
	// to dereference the pointer using the * operator to get the underlying value
	// before assigning it to our video record.
	if input.Title != nil {
		video.Title = *input.Title
	}
	// We also do the same for the other fields in the input struct.
	if input.Year != nil {
		video.Year = *input.Year
	}
	if input.Runtime != nil {
		video.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		video.Genres = input.Genres // Note that we don't need to dereference a slice.
	}
	v := validator.New()
	if data.ValidateVideo(v, video); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Videos.Update(video)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"video": video}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteVideoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the video ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the video from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	err = app.models.Videos.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "video successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
