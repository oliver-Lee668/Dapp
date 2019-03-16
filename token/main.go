package main

import (
	"token/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	router.Router(r)
	r.Run("0.0.0.0:8080")
}
