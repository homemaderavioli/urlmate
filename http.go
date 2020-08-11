package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	domainName = "localhost"
	ip         = "localhost"
	port       = 8080
)

type newURL struct {
	URL       string `json:"url" binding:"required"`
	ShortName string `json:"short_name"`
}

func createNewURL(c *gin.Context) {
	var url newURL
	if err := c.BindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rc := RiakStorageClient{
		IP:             "localhost",
		ShortURLBucket: "short_urls",
	}
	shortURL, err := SaveURL(rc, url.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"short_url": fmt.Sprintf("http://%s:%d/%s", domainName, port, shortURL)})
}

func redirectURL(c *gin.Context) {
	urlKey := c.Param("action")[1:]

	rc := RiakStorageClient{
		IP:             "localhost",
		ShortURLBucket: "short_urls",
	}
	url, err := FindURL(rc, urlKey)
	if err != nil {
		c.String(http.StatusNotFound, urlKey)
		return
	}
	c.Redirect(http.StatusFound, url)
}

func main() {
	rtr := gin.Default()
	rtr.POST("/create_url", createNewURL)
	rtr.GET("/*action", redirectURL)
	rtr.Run()
}
