package main

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"time"
)

// URL is the structure to hold information about a shorten url
type URL struct {
	OriginalURL    string
	CreationDate   string
	ExpirationDate string
}

// URLStorage interface
type URLStorage interface {
	CreateURL(key string, value []byte) (string, error)
	GetURL(key string) ([]byte, error)
}

const (
	twoYears = 8766 * 2
)

// SaveURL api
func SaveURL(sc URLStorage, originalURL string) (string, error) {
	creationDate := time.Now()
	expirationDate := creationDate.Add(time.Hour * twoYears)
	var url = &URL{
		OriginalURL:    originalURL,
		CreationDate:   creationDate.String(),
		ExpirationDate: expirationDate.String(),
	}
	jsonURL, err := json.Marshal(url)
	if err != nil {
		return "", err
	}

	urlKey := generateURLKey(6)
	key, err := sc.CreateURL(urlKey, jsonURL)
	if err != nil {
		return "", err
	}
	return key, nil
}

// FindURL api
func FindURL(sc URLStorage, urlKey string) (string, error) {
	rawURL, err := sc.GetURL(urlKey)
	if err != nil {
		return "", err
	}
	var url = &URL{}
	err = json.Unmarshal(rawURL, url)
	if err != nil {
		return "", err
	}
	return url.OriginalURL, nil
}

func generateURLKey(length int) string {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(buf)
}
