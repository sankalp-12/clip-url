package routes

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/sankalp-12/clip-url/controllers"
)

func SetupRouter(db *badger.DB, influxdb api.WriteAPIBlocking) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("api/v1")
	{
		v1.GET("resolve/:url", func(c *gin.Context) {
			start := time.Now()

			// Main handler
			controllers.ResolveURL(db, c, influxdb)

			// latency calculation
			latency := time.Since(start)

			tags := map[string]string{
				"Request": "Resolve",
			}
			fields := map[string]interface{}{
				"latency": latency,
			}
			point := write.NewPoint("latency_measurement", tags, fields, time.Now())

			if err := influxdb.WritePoint(context.Background(), point); err != nil {
				fmt.Println("[Error]: /resolve: error sending latency to influx")
			}
		})
		v1.POST("shorten", func(c *gin.Context) {
			start := time.Now()

			// Main Handler execution
			controllers.ShortenURL(db, c)

			latency := time.Since(start)

			tags := map[string]string{
				"Request": "Shorten",
			}
			fields := map[string]interface{}{
				"latency": latency,
			}
			point := write.NewPoint("latency_measurement", tags, fields, time.Now())

			if err := influxdb.WritePoint(context.Background(), point); err != nil {
				fmt.Println("[Error]: /shorten: error sending latency to influx")
			}
		})
	}

	return r
}
