package main

import (
	"net/http"

	"github.com/contemplator/HungStockKeeper/backend/database"
	"github.com/contemplator/HungStockKeeper/backend/handlers"
	"github.com/contemplator/HungStockKeeper/backend/middleware"
	"github.com/gin-gonic/gin"
)

const serverVersion = "0.1.0"

func main() {
	database.InitDB()

	r := gin.Default()

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": serverVersion,
		})
	})

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/logout", handlers.Logout)

	// Protected routes
	authorized := r.Group("/")
	authorized.Use(middleware.RequireAuth)
	{
		authorized.GET("/me", handlers.GetProfile)

		// Holdings routes
		authorized.POST("/holdings", handlers.CreateHolding)
		authorized.GET("/holdings", handlers.GetHoldings)
		authorized.GET("/holdings/:id", handlers.GetHolding)
		authorized.PUT("/holdings/:id", handlers.UpdateHolding)
		authorized.DELETE("/holdings/:id", handlers.DeleteHolding)
	}

	r.Run(":8090") // listen and serve on 0.0.0.0:8090
}
