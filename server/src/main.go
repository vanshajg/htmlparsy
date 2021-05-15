package main

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	// serving client build files
	r.Use(static.Serve("/", static.LocalFile("./web", true)))

	api := r.Group("/api")
	api.GET("/data", getData)

	r.Run()
}
