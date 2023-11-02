package controllers

import (
	"net/http"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/sankalp-12/clip-url/utils"
)

func ResolveURL(db *badger.DB, c *gin.Context) {
	url := c.Param("url")

	if flag, value := utils.CheckCollisions(db, []byte("r:"), url); flag {
		c.Redirect(http.StatusFound, value)
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error":   "URL not found",
		"message": "The requested URL does not exist in the DB",
	})
}
