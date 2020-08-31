package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleNewURL(t *testing.T) {
	srv := &server{
		domain: "localhost",
		port:   ":8080",
		db:     &StubStorage{},
		router: http.NewServeMux(),
	}
	srv.routes()

	data := `{"url":"http://www.goole.com"}`
	req := httptest.NewRequest("POST", "/create_url", strings.NewReader(data))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestHandleRedirectURL(t *testing.T) {
	expectedURL := "https://www.google.com"
	var urlData = &url{
		OriginalURL:    expectedURL,
		CreationDate:   "",
		ExpirationDate: "",
	}
	data, _ := json.Marshal(urlData)
	sc := StubStorage{
		url: data,
	}
	srv := &server{
		domain: "localhost",
		port:   ":8080",
		db:     sc,
		router: http.NewServeMux(),
	}
	srv.routes()

	req := httptest.NewRequest("GET", "/RANDOM_STRING", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusFound {
		t.Errorf("expected status code %d, got %d", http.StatusFound, w.Code)
	}
}
