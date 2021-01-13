package main

func (s *server) routes() {
	s.router.HandleFunc("/", s.handlerIndex()).Methods("GET")
	s.router.HandleFunc("/api/token", s.handlerTokenCreate()).Methods("POST")
	s.router.HandleFunc("/api/movie/{id:[0-9]+}", s.loggedOnly(s.handleMovieDetail())).Methods("GET")
	s.router.HandleFunc("/api/movies", s.loggedOnly(s.handleMovieList())).Methods("GET")
	s.router.HandleFunc("/api/movies", s.loggedOnly(s.handleMovieCreate())).Methods("POST")
}
