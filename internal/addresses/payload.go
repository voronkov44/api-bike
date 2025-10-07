package addresses

type AddressCreateRequest struct {
	Label     string `json:"label" validate:"required" example:"home"`
	Apartment string `json:"apartment,omitempty" example:"277"`
	Floor     string `json:"floor,omitempty" example:"3"`
	Entrance  string `json:"entrance,omitempty" example:"1"`
	Street    string `json:"street" validate:"required" example:"Lenina 10"`
	City      string `json:"city" validate:"required" example:"Moscow"`
	Phone     string `json:"phone,omitempty" validate:"required" example:"+7 800 555 35 55"`
	Comment   string `json:"comment,omitempty" example:"Легче позвонить, чем у кого то занимать"`
}

type AddressUpdateRequest struct {
	Label     *string `json:"label,omitempty" example:"work"`
	Apartment *string `json:"apartment,omitempty" example:"33"`
	Floor     *string `json:"floor,omitempty" example:"27"`
	Entrance  *string `json:"entrance,omitempty" example:"A"`
	Street    *string `json:"street,omitempty" example:"Pushkina 5"`
	City      *string `json:"city,omitempty" example:"Saint-Petersburg"`
	Phone     *string `json:"phone,omitempty" example:"+7 952 812 52 52"`
	Comment   *string `json:"comment,omitempty" example:"Не работает домофон"`
}

type AddressResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Label     string `json:"label"`
	Apartment string `json:"apartment,omitempty"`
	Floor     string `json:"floor,omitempty"`
	Entrance  string `json:"entrance,omitempty"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Phone     string `json:"phone,omitempty"`
	Comment   string `json:"comment,omitempty"`
	CreatedAt string `json:"created_at"`
}

// AdminAddressesResponse — формат ответа для админа
type AdminAddressesResponse struct {
	Addresses  []AddressResponse `json:"addresses"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}
