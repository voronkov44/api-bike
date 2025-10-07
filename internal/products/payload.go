package products

type ProductCreateRequest struct {
	Name        string   `json:"name" validate:"required,min=1" example:"Маргарита"`
	Type        string   `json:"type" validate:"omitempty,max=64" example:"pizza"`
	Tags        []string `json:"tags" validate:"omitempty,dive,required" example:"[\"italian\",\"popular\"]"`
	Price       int      `json:"price" validate:"required,gt=0" example:"499"`
	Ingredients []string `json:"ingredients" example:"[\"томатный соус\",\"моцарелла\",\"помидоры\",\"базилик\"]"`
	Image       string   `json:"image" validate:"omitempty,url" example:"https://example.com/image.jpg"`
	Rating      float64  `json:"rating" validate:"gte=0,lte=5" example:"4.5"`
}

type ProductUpdateRequest struct {
	Name        *string   `json:"name" validate:"omitempty,min=1"`
	Type        *string   `json:"type" validate:"omitempty,max=64"`
	Tags        *[]string `json:"tags" validate:"omitempty,dive,required"`
	Price       *int      `json:"price" validate:"omitempty,gt=0"`
	Ingredients *[]string `json:"ingredients"`
	Image       *string   `json:"image" validate:"omitempty,url"`
	Rating      *float64  `json:"rating" validate:"omitempty,gte=0,lte=5"`
}

type ProductSlugUpdateRequest struct {
	Slug string `json:"slug" validate:"required,min=1" example:"margarita-2025"`
}
