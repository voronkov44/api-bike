package products

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string   `json:"name"`
	Price       int      `json:"price"`
	Ingredients []string `json:"ingredients" gorm:"type:text[]"`
	Image       string   `json:"image"`
	Rating      float64  `json:"rating"`
}
