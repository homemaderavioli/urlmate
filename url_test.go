package main

import (
	"encoding/json"
	"testing"
)

func TestSaveURL(t *testing.T) {
	expectedKey := "dGVzdA"
	sc := StubStorage{
		Key: expectedKey,
	}
	key, err := SaveURL(sc, "https://www.google.com")
	if err != nil {
		t.Errorf("%s", err)
	}
	if key != expectedKey {
		t.Errorf("expected %s, got %s", expectedKey, key)
	}
}

func TestFindURL(t *testing.T) {
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

	url, err := FindURL(sc, "dGVzdA")
	if err != nil {
		t.Errorf("%s", err)
	}
	if url != expectedURL {
		t.Errorf("expected %s, got %s", expectedURL, url)
	}
}
