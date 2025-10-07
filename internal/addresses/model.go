package addresses

import "gorm.io/gorm"

type Address struct {
	gorm.Model `swaggerignore:"true"`
	UserID     uint   `json:"user_id" gorm:"index"`
	Label      string `json:"label"` // home-work
	Apartment  string `json:"apartment"`
	Floor      string `json:"floor"`
	Entrance   string `json:"entrance"`
	Street     string `json:"street"`
	City       string `json:"city"`
	Phone      string `json:"phone"`
	Comment    string `json:"comment"`
}
