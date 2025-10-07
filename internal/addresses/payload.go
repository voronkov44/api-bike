package addresses

type AddressCreateRequest struct {
	Label     string `json:"label" validate:"required"`
	Apartment string `json:"apartment,omitempty"`
	Floor     string `json:"floor,omitempty"`
	Entrance  string `json:"entrance,omitempty"`
	Street    string `json:"street" validate:"required"`
	City      string `json:"city" validate:"required"`
	Phone     string `json:"phone,omitempty" validate:"required"`
	Comment   string `json:"comment,omitempty"`
}

type AddressUpdateRequest struct {
	Label     *string `json:"label,omitempty"`
	Apartment *string `json:"apartment,omitempty"`
	Floor     *string `json:"floor,omitempty"`
	Entrance  *string `json:"entrance,omitempty"`
	Street    *string `json:"street,omitempty"`
	City      *string `json:"city,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Comment   *string `json:"comment,omitempty"`
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
