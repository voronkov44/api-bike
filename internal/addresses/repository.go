package addresses

import (
	"bike/pkg/db"
	"time"
)

type AddressRepository struct {
	database *db.Db
}

func NewAddressRepository(database *db.Db) *AddressRepository {
	return &AddressRepository{database: database}
}

func (r *AddressRepository) Create(a *Address) (*Address, error) {
	result := r.database.DB.Create(a)
	if result.Error != nil {
		return nil, result.Error
	}
	return a, nil
}

func (r *AddressRepository) ListByUserID(userID uint) ([]Address, error) {
	var list []Address
	result := r.database.DB.
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

func (r *AddressRepository) FindByID(id uint) (*Address, error) {
	var a Address
	result := r.database.DB.First(&a, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &a, nil
}

func (r *AddressRepository) Update(a *Address) (*Address, error) {
	result := r.database.DB.Save(a)
	if result.Error != nil {
		return nil, result.Error
	}
	return a, nil
}

func (r *AddressRepository) DeleteByID(id uint) error {
	result := r.database.DB.Delete(&Address{}, id)
	return result.Error
}

// ToResponse helpers func
func ToResponse(a *Address) AddressResponse {
	created := ""
	if !a.CreatedAt.IsZero() {
		created = a.CreatedAt.Format(time.RFC3339)
	}
	return AddressResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		Label:     a.Label,
		Apartment: a.Apartment,
		Floor:     a.Floor,
		Entrance:  a.Entrance,
		Street:    a.Street,
		City:      a.City,
		Phone:     a.Phone,
		Comment:   a.Comment,
		CreatedAt: created,
	}
}
