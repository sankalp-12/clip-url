package controllers

import (
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
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Request might contain missing fields",
		})
		return
	}

	if !govalidator.IsURL(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Invalid URL",
		})
		return
	}

	if !utils.RemoveDomainError(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "URL requested is forbidden",
		})
		return
	}
	body.URL = utils.EnforceHTTP(body.URL)

	var newURL string
	if body.Custom != "" {
		if flag, _ := utils.CheckCollisions(db, []byte("r:"), body.Custom); !flag {
			newURL = body.Custom
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad request",
				"message": "Custom URL is already in use",
			})
			return
		}
	} else {
		if flag, oldURL := utils.CheckCollisions(db, []byte("w:"), body.URL); flag {
			c.JSON(http.StatusOK, Response{
				URL:    body.URL,
				NewURL: oldURL,
			})
			return
		}

		newURL = uuid.New().String()[:6]
		flag, _ := utils.CheckCollisions(db, []byte("r:"), newURL)
		for flag {
			newURL = uuid.New().String()[:6]
			flag, _ = utils.CheckCollisions(db, []byte("r:"), newURL)
		}
	}

	if flag, message := utils.WriteToDB(db, body.URL, newURL); !flag {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		URL:    body.URL,
		NewURL: newURL,
	})
}
