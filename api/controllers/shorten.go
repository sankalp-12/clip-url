package controllers

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sankalp-12/clip-url/models"
	"github.com/sankalp-12/clip-url/utils"
)

func ShortenURL(db *badger.DB, c *gin.Context) {
	body := new(models.Request)
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			URL:     body.URL,
			NewURL:  "",
			Message: "Bad request: Contains missing fields",
		})
		return
	}

	if !govalidator.IsURL(body.URL) {
		c.JSON(http.StatusBadRequest, models.Response{
			URL:     body.URL,
			NewURL:  "",
			Message: "Bad request: URL is not valid",
		})
		return
	}

	if !utils.RemoveDomainError(body.URL) {
		c.JSON(http.StatusBadRequest, models.Response{
			URL:     body.URL,
			NewURL:  "",
			Message: "Bad request: URL is forbidden",
		})
		return
	}
	body.URL = utils.EnforceHTTP(body.URL)

	var newURL string
	if body.Custom != "" {
		if flag, _ := utils.CheckCollisions(db, []byte("r:"), body.Custom); !flag {
			newURL = body.Custom
		} else {
			c.JSON(http.StatusBadRequest, models.Response{
				URL:     body.URL,
				NewURL:  "",
				Message: "Bad request: Custom URL already exists",
			})
			return
		}
	} else {
		if flag, oldURL := utils.CheckCollisions(db, []byte("w:"), body.URL); flag {
			c.JSON(http.StatusOK, models.Response{
				URL:     body.URL,
				NewURL:  oldURL,
				Message: "Success: Requested URL already exists, returning old URL",
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
		c.JSON(http.StatusBadRequest, models.Response{
			URL:     body.URL,
			NewURL:  "",
			Message: "Internal server error: " + message,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		URL:     body.URL,
		NewURL:  newURL,
		Message: "Success: URL shortened",
	})
}
