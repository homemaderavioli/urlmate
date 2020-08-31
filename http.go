package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type server struct {
	domain string
	port   string
	db     Storage
	router *http.ServeMux
}

func decode(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return err
	}
	if valid, ok := v.(interface {
		OK() error
	}); ok {
		err := valid.OK()
		if err != nil {
			return err
		}
	}
	return nil
}

func respondErr(w http.ResponseWriter, r *http.Request, err error, code int) {
	errObj := struct {
		Error string `json:"error"`
	}{Error: err.Error()}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	err = json.NewEncoder(w).Encode(errObj)
	if err != nil {
		fmt.Errorf("repondErr: %s", err)
	}
}

func respond(w http.ResponseWriter, r *http.Request, v interface{}, code int) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(v)
	if err != nil {
		respondErr(w, r, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Errorf("respond: %s", err)
	}
}

func (s *server) handleNewURL() http.HandlerFunc {
	var request struct {
		URL       string `json:"url"`
		ShortName string `json:"short_name"`
	}
	var response struct {
		ShortURL string `json:"short_url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		err := decode(r, &request)
		if err != nil {
			respondErr(w, r, err, http.StatusBadRequest)
		}
		shortURL, err := SaveURL(s.db, request.URL)
		if err != nil {
			respondErr(w, r, err, http.StatusBadRequest)
		}
		response.ShortURL = fmt.Sprintf("http://%s%s/%s", s.domain, s.port, shortURL)
		respond(w, r, response, http.StatusCreated)
	}
}

func (s *server) handleRedirectURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := FindURL(s.db, r.URL.Path[1:])
		if err != nil {
			respondErr(w, r, err, http.StatusNotFound)
		}
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
