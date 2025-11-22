package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const serverVersion = "0.1.0"

func main() {
	r := gin.Default()

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"version": serverVersion,
		})
	})

	r.Run(":8090") // listen and serve on 0.0.0.0:8090
}
