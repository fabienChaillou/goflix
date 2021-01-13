package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStore struct {
	movieId int64
	movies  []*Movie
}

func (t testStore) Open() error {
	return nil
}

func (t testStore) Close() error {
	return nil
}

func (t testStore) GetMovies() ([]*Movie, error) {
	return t.movies, nil
}

func (t testStore) GetMovieById(id int64) (*Movie, error) {
	for _, m := range t.movies {
		if m.ID == id {
			return m, nil
		}
	}
	return nil, nil
}

func (t testStore) CreateMovie(m *Movie) error {
	t.movieId++
	m.ID = t.movieId
	t.movies = append(t.movies, m)
	return nil
}

func TestMovieCreateUnit(t *testing.T) {
	srv := newServer()
	srv.store = &testStore{}

	p := struct {
		Title       string `json:"title"`
		ReleaseDate string `json:"release_date"`
		Duration    int    `json:"duration"`
		TrailerURL  string `json:"trailer_url"`
	}{
		Title:       "inception",
		ReleaseDate: "2010-01-12",
		Duration:    125,
		TrailerURL:  "http://test.io",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(p)
	assert.Nil(t, err)

	r := httptest.NewRequest("POST", "/api/movies", &buf)
	w := httptest.NewRecorder()

	srv.handleMovieCreate()(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestMovieCreateIntegration(t *testing.T) {
	srv := newServer()
	srv.store = &testStore{}

	p := struct {
		Title       string `json:"title"`
		ReleaseDate string `json:"release_date"`
		Duration    int    `json:"duration"`
		TrailerURL  string `json:"trailer_url"`
	}{
		Title:       "inception",
		ReleaseDate: "2010-01-12",
		Duration:    125,
		TrailerURL:  "http://test.io",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(p)
	assert.Nil(t, err)

	r := httptest.NewRequest("POST", "/api/movies", &buf)
	w := httptest.NewRecorder()

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTExOTcxNjUsImlhdCI6MTU5MTE5MzU2NSwidXNlcm5hbWUiOiJyb2JpbiJ9.QnGWAATeo1ey8CVBwngazJF_dkGHSvJVf6tCGTHHwC4"
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	srv.serveHTTP(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
}
