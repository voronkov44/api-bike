package products

type ProductCreateRequest struct {
	Name        string   `json:"name" validate:"required,min=1"`
	Price       int      `json:"price" validate:"required,gt=0"`
	Ingredients []string `json:"ingredients"`
	Image       string   `json:"image" validate:"omitempty,url"`
	Rating      float64  `json:"rating" validate:"gte=0,lte=5"`
}
