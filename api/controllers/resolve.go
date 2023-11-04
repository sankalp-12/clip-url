package controllers

import (
	"net/http"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/sankalp-12/clip-url/models"
	"github.com/sankalp-12/clip-url/utils"
)

func ResolveURL(db *badger.DB, c *gin.Context) {
	url := c.Param("url")

	if flag, value := utils.CheckCollisions(db, []byte("r:"), url); flag {
		c.Redirect(http.StatusFound, value)
		return
	}

	c.JSON(http.StatusBadRequest, models.Response{
		URL:     url,
		NewURL:  "",
		Message: "Bad request: The URL does not exist",
	})
}
