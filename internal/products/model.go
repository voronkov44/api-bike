package products

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model  `swaggerignore:"true"`
	Slug        string         `json:"slug" gorm:"size:128;uniqueIndex;not null"`
	Name        string         `json:"name" gorm:"not null;uniqueIndex"`
	Type        string         `json:"type" gorm:"size:64;index"`
	Price       int            `json:"price"`
	Ingredients pq.StringArray `json:"ingredients" gorm:"type:text[]" swaggerignore:"true"`
	Tags        pq.StringArray `json:"tags" gorm:"type:text[]" swaggerignore:"true"`
	Image       string         `json:"image"`
	Rating      float64        `json:"rating"`
}
