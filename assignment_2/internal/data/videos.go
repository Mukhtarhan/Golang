package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"assignment_2.alexedwards.net/internal/validator"
	"github.com/lib/pq"
)

type Video struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func (m VideoModel) Insert(video *Video) error {
	// Define the SQL query for inserting a new record in the videos table and returning
	// the system-generated data.
	query := `
INSERT INTO videos (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version`
	// Create an args slice containing the values for the placeholder parameters from
	// the video struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []interface{}{video.Title, video.Year, video.Runtime, pq.Array(video.Genres)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&video.ID, &video.CreatedAt, &video.Version)
}

func ValidateVideo(v *validator.Validator, video *Video) {
	v.Check(video.Title != "", "title", "must be provided")
	v.Check(len(video.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(video.Year != 0, "year", "must be provided")
	v.Check(video.Year >= 1888, "year", "must be greater than 1888")
	v.Check(video.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(video.Runtime != 0, "runtime", "must be provided")
	v.Check(video.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(video.Genres != nil, "genres", "must be provided")
	v.Check(len(video.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(video.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(video.Genres), "genres", "must not contain duplicate values")
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

type VideoModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the Videos table.

// Add a placeholder method for fetching a specific record from the Videos table.
func (m VideoModel) Get(id int64) (*Video, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the video data.
	query := `
SELECT id, created_at, title, year, runtime, genres, version
FROM videos
WHERE id = $1`
	// Declare a video struct to hold the data returned by the query.
	var video Video

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// video struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&video.ID,
		&video.CreatedAt,
		&video.Title,
		&video.Year,
		&video.Runtime,
		pq.Array(&video.Genres),
		&video.Version,
	)
	// Handle any errors. If there was no matching video found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the video struct.
	return &video, nil

}

// Add a placeholder method for updating a specific record in the Videos table.
func (m VideoModel) Update(Video *Video) error {
	query := `
UPDATE Videos
SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		Video.Title,
		Video.Year,
		Video.Runtime,
		pq.Array(Video.Genres),
		Video.ID,
		Video.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the Video struct.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&Video.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil

}
func (m VideoModel) GetAll(title string, genres []string, filters Filters) ([]*Video, error) {
	// Construct the SQL query to retrieve all Video records.
	query := `
SELECT id, created_at, title, year, runtime, genres, version
FROM Videos
ORDER BY id`
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	// Initialize an empty slice to hold the Video data.
	Videos := []*Video{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Video struct to hold the data for an individual Video.
		var Video Video
		// Scan the values from the row into the Video struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&Video.ID,
			&Video.CreatedAt,
			&Video.Title,
			&Video.Year,
			&Video.Runtime,
			pq.Array(&Video.Genres),
			&Video.Version,
		)
		if err != nil {
			return nil, err
		}

		// Add the Video struct to the slice.
		Videos = append(Videos, &Video)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK, then return the slice of Videos.
	return Videos, nil
}

// Add a placeholder method for deleting a specific record from the Videos table.
func (m VideoModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
DELETE FROM videos
WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the Videos table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}
