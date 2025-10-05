package products

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name"`
	Price       int            `json:"price"`
	Ingredients pq.StringArray `json:"ingredients" gorm:"type:text[]"`
	Image       string         `json:"image"`
	Rating      float64        `json:"rating"`
}
