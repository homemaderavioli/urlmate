package main

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"time"
)

// Storage is the interface url data
type Storage interface {
	CreateURL(key string, value []byte) (string, error)
	GetURL(key string) ([]byte, error)
}

// url is the structure to hold information about a shorten url
type url struct {
	OriginalURL    string
	CreationDate   string
	ExpirationDate string
}

const (
	twoYears = 8766 * 2
)

// SaveURL api
func SaveURL(sc Storage, originalURL string) (string, error) {
	creationDate := time.Now()
	expirationDate := creationDate.Add(time.Hour * twoYears)
	var url = &url{
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
func FindURL(sc Storage, urlKey string) (string, error) {
	rawURL, err := sc.GetURL(urlKey)
	if err != nil {
		return "", err
	}
	var url = &url{}
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
