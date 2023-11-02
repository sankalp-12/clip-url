package controllers

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sankalp-12/clip-url/utils"
)

type Request struct {
	URL    string `json:"url"`
	Custom string `json:"custom,omitempty"`
}

type Response struct {
	URL    string `json:"url"`
	NewURL string `json:"new_url"`
}

func ShortenURL(db *badger.DB, c *gin.Context) {
	body := new(Request)
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Request might contain missing fields",
		})
		return
	}

	if !govalidator.IsURL(body.URL) || !utils.RemoveDomainError(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Invalid URL",
		})
		return
	}
	body.URL = utils.EnforceHTTP(body.URL)

	var newURL string

	if flag, oldURL := utils.CheckCollisions(db, []byte("w:"), body.URL); flag {
		c.JSON(http.StatusOK, Response{
			URL:    body.URL,
			NewURL: oldURL,
		})
		return
	}

	if body.Custom != "" {
		if flag, _ := utils.CheckCollisions(db, []byte("r:"), body.Custom); flag {
			newURL = body.Custom
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad request",
				"message": "Custom URL is already in use",
			})
			return
		}
	} else {
		newURL = uuid.New().String()[:6]
		flag, _ := utils.CheckCollisions(db, []byte("r:"), newURL)
		for flag {
			newURL = uuid.New().String()[:6]
			flag, _ = utils.CheckCollisions(db, []byte("r:"), newURL)
		}
	}

	txn := db.NewTransaction(true)
	defer txn.Discard()

	key := []byte(fmt.Sprintf("w:%s", body.URL))
	value := []byte(newURL)
	err := txn.Set(key, value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Unable to write to the DB",
		})
		return
	}

	key = []byte(fmt.Sprintf("r:%s", newURL))
	value = []byte(body.URL)
	err = txn.Set(key, value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Unable to write to the DB",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		URL:    body.URL,
		NewURL: newURL,
	})
}
