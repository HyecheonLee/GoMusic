package models

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model
	Image       string  `json:"img"`
	SmallImage  string  `gorm:"column:smallimg" json:"small_img"`
	ImagAlt     string  `json:"imgalt" gorm:"column:imgalt"`
	Price       float64 `json:"price"`
	Promotion   float64 `json:"promotion"` //sql.NullFloat64
	PoructName  string  `gorm:"column:productname" json:"productname"`
	Description string
}
