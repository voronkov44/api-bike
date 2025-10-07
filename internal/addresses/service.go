package addresses

import (
	"context"
	"errors"
	"math"

	"bike/internal/users"
)

var (
	ErrAddressNotFound = errors.New("address not found")
	ErrForbidden       = errors.New("forbidden")
)

type AddressService struct {
	repo     *AddressRepository
	userRepo *users.UserRepository
}

func NewAddressService(repo *AddressRepository, userRepo *users.UserRepository) *AddressService {
	return &AddressService{
		repo:     repo,
		userRepo: userRepo,
	}
}

// CreateAddress создаёт адрес для пользователя, найденного по email.
func (s *AddressService) CreateAddress(ctx context.Context, userEmail string, in AddressCreateRequest) (*Address, error) {
	user, err := s.userRepo.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}
	a := &Address{
		UserID:    user.ID,
		Label:     in.Label,
		Apartment: in.Apartment,
		Floor:     in.Floor,
		Entrance:  in.Entrance,
		Street:    in.Street,
		City:      in.City,
		Phone:     in.Phone,
		Comment:   in.Comment,
	}
	return s.repo.Create(a)
}

// ListAddress возвращает все адреса конкретного пользователя.
func (s *AddressService) ListAddress(ctx context.Context, userEmail string) ([]Address, error) {
	user, err := s.userRepo.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}
	return s.repo.ListByUserID(user.ID)
}

// UpdateAddress обновляет поля адреса, только если он принадлежит пользователю.
func (s *AddressService) UpdateAddress(ctx context.Context, userEmail string, id uint, in AddressUpdateRequest) (*Address, error) {
	user, err := s.userRepo.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}
	addr, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrAddressNotFound
	}
	if addr.UserID != user.ID {
		return nil, ErrForbidden
	}

	// Накатываем только те поля, которые пришли (nil — пропустить)
	if in.Label != nil {
		addr.Label = *in.Label
	}
	if in.Apartment != nil {
		addr.Apartment = *in.Apartment
	}
	if in.Floor != nil {
		addr.Floor = *in.Floor
	}
	if in.Entrance != nil {
		addr.Entrance = *in.Entrance
	}
	if in.Street != nil {
		addr.Street = *in.Street
	}
	if in.City != nil {
		addr.City = *in.City
	}
	if in.Phone != nil {
		addr.Phone = *in.Phone
	}
	if in.Comment != nil {
		addr.Comment = *in.Comment
	}

	return s.repo.Update(addr)
}

// DeleteAddress удаляет адрес только если он принадлежит пользователю.
func (s *AddressService) DeleteAddress(ctx context.Context, userEmail string, id uint) error {
	user, err := s.userRepo.FindByEmail(userEmail)
	if err != nil {
		return err
	}
	addr, err := s.repo.FindByID(id)
	if err != nil {
		return ErrAddressNotFound
	}
	if addr.UserID != user.ID {
		return ErrForbidden
	}
	return s.repo.DeleteByID(id)
}

func (s *AddressService) ListAllAdmin(ctx context.Context, userID uint, city, street, phone, label string, page, limit int) (items []Address, total int64, totalPages int, err error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	items, total, err = s.repo.ListAll(userID, city, street, phone, label, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	if total == 0 {
		totalPages = 0
	} else {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}
	return items, total, totalPages, nil
}
