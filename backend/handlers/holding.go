package handlers

import (
	"net/http"
	"time"

	"github.com/contemplator/HungStockKeeper/backend/database"
	"github.com/contemplator/HungStockKeeper/backend/models"
	"github.com/gin-gonic/gin"
)

// CreateHolding adds a new holding for the authenticated user
func CreateHolding(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(models.User)

	var input models.CreateHoldingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default purchase date to now if not provided
	if input.PurchaseDate.IsZero() {
		input.PurchaseDate = time.Now()
	}

	holding := models.Holding{
		UserID:       currentUser.ID,
		Symbol:       input.Symbol,
		Quantity:     input.Quantity,
		CostBasis:    input.CostBasis,
		PurchaseDate: input.PurchaseDate,
		BrokerageId:  input.BrokerageId,
		Note:         input.Note,
	}

	if err := database.DB.Create(&holding).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create holding"})
		return
	}

	c.JSON(http.StatusCreated, holding)
}

// GetHoldings returns all holdings for the authenticated user
func GetHoldings(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(models.User)

	var holdings []models.Holding
	if err := database.DB.Where("user_id = ?", currentUser.ID).Find(&holdings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch holdings"})
		return
	}

	c.JSON(http.StatusOK, holdings)
}

// GetHolding returns a single holding by ID
func GetHolding(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(models.User)

	var holding models.Holding
	if err := database.DB.Where("id = ? AND user_id = ?", c.Param("id"), currentUser.ID).First(&holding).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Holding not found"})
		return
	}

	c.JSON(http.StatusOK, holding)
}

// UpdateHolding updates an existing holding
func UpdateHolding(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(models.User)

	var holding models.Holding
	if err := database.DB.Where("id = ? AND user_id = ?", c.Param("id"), currentUser.ID).First(&holding).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Holding not found"})
		return
	}

	var input models.UpdateHoldingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if input.Symbol != "" {
		holding.Symbol = input.Symbol
	}
	if input.Quantity != 0 {
		holding.Quantity = input.Quantity
	}
	if input.CostBasis != 0 {
		holding.CostBasis = input.CostBasis
	}
	if !input.PurchaseDate.IsZero() {
		holding.PurchaseDate = input.PurchaseDate
	}
	if input.BrokerageId != nil {
		holding.BrokerageId = input.BrokerageId
	}
	holding.Note = input.Note

	if err := database.DB.Save(&holding).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update holding"})
		return
	}

	c.JSON(http.StatusOK, holding)
}

// DeleteHolding removes a holding
func DeleteHolding(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(models.User)

	var holding models.Holding
	if err := database.DB.Where("id = ? AND user_id = ?", c.Param("id"), currentUser.ID).First(&holding).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Holding not found"})
		return
	}

	if err := database.DB.Delete(&holding).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete holding"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Holding deleted successfully"})
}
