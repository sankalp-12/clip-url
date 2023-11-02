package routes

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/sankalp-12/clip-url/controllers"
)

func SetupRouter(db *badger.DB) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("api/v1")
	{
		v1.GET("resolve/:url", func(c *gin.Context) {
			controllers.ResolveURL(db, c)
		})
		v1.POST("shorten", func(c *gin.Context) {
			controllers.ShortenURL(db, c)
		})
	}

	return r
}
