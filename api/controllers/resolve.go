package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/sankalp-12/clip-url/models"
	"github.com/sankalp-12/clip-url/utils"
)

func ResolveURL(db *badger.DB, c *gin.Context, write_api api.WriteAPIBlocking) {
	url := c.Param("url")

	if flag, value := utils.CheckCollisions(db, []byte("r:"), url); flag {
		tags := map[string]string{
			"Resolve": "Hit",
		}
		fields := map[string]interface{}{
			"url": url,
		}
		point := write.NewPoint("Metrics", tags, fields, time.Now())

		if err := write_api.WritePoint(context.Background(), point); err != nil {
			fmt.Println("[Error]: /resolve: error sending timestamp of resolve hit to influx")
		}
		c.Redirect(http.StatusFound, value)
		return
	}

	c.JSON(http.StatusBadRequest, models.Response{
		URL:     url,
		NewURL:  "",
		Message: "Bad request: The URL does not exist",
	})

}
