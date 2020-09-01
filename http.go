package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
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
	fmt.Errorf("respond error: %v", err)
	errObj := struct {
		Error string `json:"error"`
	}{Error: err.Error()}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	err = json.NewEncoder(w).Encode(errObj)
	if err != nil {
		fmt.Errorf("respond err: %s", err)
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
		if r.Method != http.MethodPost {
			respondErr(w, r, errors.New("bad request"), http.StatusBadRequest)
			return
		}
		err := decode(r, &request)
		if err != nil {
			respondErr(w, r, err, http.StatusBadRequest)
			return
		}
		if request.URL == "" {
			respondErr(w, r, errors.New("field url is required"), http.StatusBadRequest)
			return
		}
		shortURL, err := SaveURL(s.db, request.URL)
		if err != nil {
			respondErr(w, r, err, http.StatusBadRequest)
			return
		}
		response.ShortURL = fmt.Sprintf("http://%s%s/%s", s.domain, s.port, shortURL)
		respond(w, r, response, http.StatusCreated)
	}
}

func (s *server) handleRedirectURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondErr(w, r, errors.New("bad request"), http.StatusBadRequest)
			return
		}
		url, err := FindURL(s.db, html.EscapeString(r.URL.Path[1:]))
		if err != nil {
			respondErr(w, r, errors.New("not found"), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
