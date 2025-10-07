package addresses

import (
	"bike/pkg/db"
	"gorm.io/gorm"
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

func (r *AddressRepository) ListAll(userID uint, city, street, phone, label string, limit, offset int) (items []Address, total int64, err error) {
	dbq := r.database.DB.Model(&Address{})

	if userID != 0 {
		dbq = dbq.Where("user_id = ?", userID)
	}
	if city != "" {
		dbq = dbq.Where("city ILIKE ?", "%"+city+"%")
	}
	if street != "" {
		dbq = dbq.Where("street ILIKE ?", "%"+street+"%")
	}
	if phone != "" {
		dbq = dbq.Where("phone ILIKE ?", "%"+phone+"%")
	}
	if label != "" {
		dbq = dbq.Where("label = ?", label)
	}

	// count total
	if err = dbq.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// fetch with pagination
	result := dbq.Order("created_at desc").Limit(limit).Offset(offset).Find(&items)
	if result.Error != nil {
		// если нет записей — возвращаем пустой слайс
		if result.Error == gorm.ErrRecordNotFound {
			return []Address{}, total, nil
		}
		return nil, 0, result.Error
	}
	return items, total, nil
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
