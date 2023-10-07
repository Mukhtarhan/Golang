package data

import (
	"encoding/json"
	"fmt"
	"time"
)

type Video struct {
	ID        int64
	CreatedAt time.Time
	Title     string
	Year      int32
	Runtime   int32
	Genres    []string
	Version   int32
}

// Implement a MarshalJSON() method on the Video struct, so that it satisfies the
// json.Marshaler interface.
func (m Video) MarshalJSON() ([]byte, error) {
	// Declare a variable to hold the custom runtime string (this will be the empty
	// string "" by default).
	var runtime string
	// If the value of the Runtime field is not zero, set the runtime variable to be a
	// string in the format "<runtime> mins".
	if m.Runtime != 0 {
		runtime = fmt.Sprintf("%d mins", m.Runtime)
	}
	// Create an anonymous struct to hold the data for JSON encoding. This has exactly
	// the same fields, types and tags as our Video struct, except that the Runtime
	// field here is a string, instead of an int32. Also notice that we don't include
	// a CreatedAt field at all (there's no point including one, because we don't want
	// it to appear in the JSON output).
	aux := struct {
		ID      int64    `json:"id"`
		Title   string   `json:"title"`
		Year    int32    `json:"year,omitempty"`
		Runtime string   `json:"runtime,omitempty"` // This is a string.
		Genres  []string `json:"genres,omitempty"`
		Version int32    `json:"version"`
	}{
		// Set the values for the anonymous struct.
		ID:      m.ID,
		Title:   m.Title,
		Year:    m.Year,
		Runtime: runtime, // Note that we assign the value from the runtime variable here.
		Genres:  m.Genres,
		Version: m.Version,
	}
	// Encode the anonymous struct to JSON, and return it.
	return json.Marshal(aux)
}
