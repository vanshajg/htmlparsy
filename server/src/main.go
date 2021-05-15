package main

import (
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

func main() {

	r := gin.Default()
	// serving client build files
	r.Use(static.Serve("/", static.LocalFile("./web", true)))
	apiCache := cache.New(5*time.Minute, 10*time.Minute)
	api := r.Group("/api")
	api.GET("/data", CacheCheck(apiCache), getData)

	r.Run()
}
