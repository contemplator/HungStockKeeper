package handlers

import (
	"net/http"

	"github.com/contemplator/HungStockKeeper/backend/database"
	"github.com/contemplator/HungStockKeeper/backend/models"
	"github.com/gin-gonic/gin"
)

// GetBrokerages returns all available brokerages
func GetBrokerages(c *gin.Context) {
	var brokerages []models.Brokerage
	if err := database.DB.Find(&brokerages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch brokerages"})
		return
	}

	c.JSON(http.StatusOK, brokerages)
}
