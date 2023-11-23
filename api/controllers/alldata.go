package controllers

import (
	"net/http"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/sankalp-12/clip-url/models"
	"github.com/sankalp-12/clip-url/utils"
)

func AllShortenData(db *badger.DB, g *gin.Context) {
	conf, data := utils.ReturnAllData(db)
	if !conf {
		g.JSON(http.StatusInternalServerError, models.ResponseData{
			Error:   1,
			Message: "Issue with BadgerDB",
		})
	}

	g.JSON(http.StatusOK, models.ResponseData{
		Error:   0,
		Message: "OK",
		Data:    data,
	})
}
