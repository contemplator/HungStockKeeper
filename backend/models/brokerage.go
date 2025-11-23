package models

type Brokerage struct {
	ID   int64  `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null;unique"`
}
