package models

import (
	"time"

	"gorm.io/gorm"
)

type Holding struct {
	gorm.Model
	UserID int64 `json:"user_id" gorm:"not null;index"`
	// 使用 gorm:"column:stock_code" 明確指定對應的資料庫欄位
	Symbol       string    `json:"symbol" binding:"required" gorm:"column:stock_code"`
	StockName    string    `json:"stock_name" gorm:"column:stock_name"`
	PurchaseDate time.Time `json:"purchase_date" gorm:"column:buy_date"`
	CostBasis    float64   `json:"cost_basis" gorm:"not null;column:buy_price"`
	Quantity     float64   `json:"quantity" gorm:"not null;column:quantity"`
	BrokerageId  *int64    `json:"brokerage_id" gorm:"column:brokerage_id"`
	Note         string    `json:"note" gorm:"size:50"`
}

type CreateHoldingInput struct {
	Symbol       string    `json:"symbol" binding:"required"`
	StockName    string    `json:"stock_name"`
	Quantity     float64   `json:"quantity" binding:"required,gt=0"`
	CostBasis    float64   `json:"cost_basis" binding:"required,gt=0"`
	PurchaseDate time.Time `json:"purchase_date"`
	BrokerageId  *int64    `json:"brokerage_id"`
	Note         string    `json:"note" binding:"max=50"`
}

type UpdateHoldingInput struct {
	Symbol       string    `json:"symbol"`
	Quantity     float64   `json:"quantity"`
	CostBasis    float64   `json:"cost_basis"`
	PurchaseDate time.Time `json:"purchase_date"`
	BrokerageId  *int64    `json:"brokerage_id"`
	Note         string    `json:"note" binding:"max=50"`
}
