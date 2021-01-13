package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type jsonMovie struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Duration    int    `json:"duration"`
	TrailerURL  string `json:"trailer_url"`
}

func (s *server) handleMovieList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		movies, err := s.store.GetMovies()
		if err != nil {
			log.Printf("Cannot load movies. err=%v\n", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		var resp = make([]jsonMovie, len(movies))
		for i, m := range movies {
			resp[i] = mapMovieToJson(m)
		}
		s.respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleMovieDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.Printf("Error:  dont parse id to int, err=%v\n", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		m, err := s.store.GetMovieById(id)
		if err != nil {
			log.Printf("Cannot load movie. err=%v\n", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		var resp = mapMovieToJson(m)
		s.respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleMovieCreate() http.HandlerFunc {
	type request struct {
		Title       string `db:"title"`
		ReleaseDate string `db:"release_date"`
		Duration    int    `db:"duration"`
		TrailerURL  string `db:"trailer_url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := s.decode(w, r, &req)
		if err != nil {
			log.Printf("Error: cannot parse movie, err=%v\n", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		// create muvie
		m := &Movie{
			ID:          0,
			Title:       req.Title,
			ReleaseDate: req.ReleaseDate,
			Duration:    req.Duration,
			TrailerURL:  req.TrailerURL,
		}

		// Persiste movie in db
		err = s.store.CreateMovie(m)
		if err != nil {
			log.Printf("Error: cannot create movie in DB, err=%v\n", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		var resq = mapMovieToJson(m)
		s.respond(w, r, resq, http.StatusCreated)

	}
}
func mapMovieToJson(m *Movie) jsonMovie {
	return jsonMovie{
		ID:          m.ID,
		Title:       m.Title,
		ReleaseDate: m.ReleaseDate,
		Duration:    m.Duration,
		TrailerURL:  m.TrailerURL,
	}
}
