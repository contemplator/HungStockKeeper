package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/contemplator/HungStockKeeper/backend/database"
	"github.com/contemplator/HungStockKeeper/backend/models"
	"github.com/gin-gonic/gin"
)

// HoldingResponse defines the structure for the API response including real-time data
type HoldingResponse struct {
	models.Holding
	CurrentPrice      float64 `json:"current_price"`       // Real-time price
	MarketValue       float64 `json:"market_value"`        // Market Value = CurrentPrice * Quantity
	ProfitLoss        float64 `json:"profit_loss"`         // Profit/Loss = (CurrentPrice - CostBasis) * Quantity
	ProfitLossPercent float64 `json:"profit_loss_percent"` // P/L %
}

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

// GetHoldings returns all holdings for the authenticated user with real-time data
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

	// 1. Collect all symbols
	var symbols []string
	uniqueSymbols := make(map[string]bool)

	for _, h := range holdings {
		// Only process Taiwan stock codes (numeric) for TWSE API
		if isTaiwanStock(h.Symbol) {
			if !uniqueSymbols[h.Symbol] {
				symbols = append(symbols, h.Symbol)
				uniqueSymbols[h.Symbol] = true
			}
		}
	}

	// 2. Batch fetch quotes from TWSE
	priceMap := make(map[string]float64)
	if len(symbols) > 0 {
		var err error
		priceMap, err = fetchStockPricesTWSE(symbols)
		if err != nil {
			// Log error internally or handle it, but don't crash the request
			fmt.Printf("Error fetching TWSE quotes: %v\n", err)
		}
	}

	// 3. Prepare response
	response := make([]HoldingResponse, len(holdings))

	// 4. Fill in data
	for i, h := range holdings {
		response[i].Holding = h

		// Get price from map (default 0 if not found)
		// The map keys are the stock codes (e.g., "2330")
		currentPrice := priceMap[h.Symbol]
		qty := h.Quantity
		cost := h.CostBasis

		// Calculations
		marketValue := currentPrice * float64(qty)
		profitLoss := (currentPrice - cost) * float64(qty)

		var profitLossPercent float64
		if cost > 0 {
			profitLossPercent = (profitLoss / (cost * float64(qty))) * 100
		}

		response[i].CurrentPrice = currentPrice
		response[i].MarketValue = marketValue
		response[i].ProfitLoss = profitLoss
		response[i].ProfitLossPercent = profitLossPercent
	}

	c.JSON(http.StatusOK, response)
}

// TwseResponse defines the structure of TWSE API response
type TwseResponse struct {
	MsgArray []struct {
		C string `json:"c"` // Code
		Z string `json:"z"` // Recent Trade Price
		Y string `json:"y"` // Yesterday Closing Price
		N string `json:"n"` // Name
	} `json:"msgArray"`
}

// fetchStockPricesTWSE fetches real-time quotes from Taiwan Stock Exchange
func fetchStockPricesTWSE(symbols []string) (map[string]float64, error) {
	if len(symbols) == 0 {
		return make(map[string]float64), nil
	}

	var queryParts []string
	for _, s := range symbols {
		// Query both TSE (Listed) and OTC (Over-the-counter) for each symbol
		// This covers most cases without needing to know the exchange type beforehand
		queryParts = append(queryParts, fmt.Sprintf("tse_%s.tw", s))
		queryParts = append(queryParts, fmt.Sprintf("otc_%s.tw", s))
	}

	apiURL := "https://mis.twse.com.tw/stock/api/getStockInfo.jsp"
	params := url.Values{}
	params.Add("ex_ch", strings.Join(queryParts, "|"))
	params.Add("json", "1")
	params.Add("delay", "0")
	params.Add("_", fmt.Sprintf("%d", time.Now().UnixMilli()))

	fullURL := apiURL + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TwseResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	prices := make(map[string]float64)
	for _, msg := range result.MsgArray {
		priceStr := msg.Z
		// If no recent trade price, fallback to yesterday's closing price
		if priceStr == "-" {
			priceStr = msg.Y
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err == nil {
			prices[msg.C] = price
		}
	}

	return prices, nil
}

// isTaiwanStock checks if the symbol is likely a Taiwan stock code (numeric)
func isTaiwanStock(symbol string) bool {
	for _, r := range symbol {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
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
