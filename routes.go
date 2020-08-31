package main

func (s *server) routes() {
	s.router.HandleFunc("/create_url", s.handleNewURL())
	s.router.HandleFunc("/", s.handleRedirectURL())
}
