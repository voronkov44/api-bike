package products

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Slug        string         `json:"slug" gorm:"size:128;uniqueIndex;not null"`
	Name        string         `json:"name" gorm:"not null;uniqueIndex"`
	Type        string         `json:"type" gorm:"size:64;index"`
	Price       int            `json:"price"`
	Ingredients pq.StringArray `json:"ingredients" gorm:"type:text[]"`
	Tags        pq.StringArray `json:"tags" gorm:"type:text[]"`
	Image       string         `json:"image"`
	Rating      float64        `json:"rating"`
}
