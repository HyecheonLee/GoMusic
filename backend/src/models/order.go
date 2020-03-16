package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Order struct {
	gorm.Model
	Product
	Customer
	CustomerID   int       `json:"customer_id" gorm:"column:customer_id"`
	ProductID    int       `json:"product_id" gorm:"column:product_id"`
	Price        float64   `gorm:"column:price" json:"sell_price"`
	PurchaseDate time.Time `gorm:"column:purchase_date" json:"purchase_date"`
}
