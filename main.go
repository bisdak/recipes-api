package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Recipe struct {
	Name        string    `json:"name"`
	Tags        []string  `json:"ingredients"`
	Ingredients []string  `json:"instructions"`
	PublishedAt time.Time `json:"publishedAt"`
}

func main() {
	router := gin.Default()
	router.Run()
}
