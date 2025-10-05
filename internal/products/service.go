package products

import (
	"context"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

var (
	ErrValidation = errors.New("validation error")
)

type ProductService interface {
	CreateProduct(ctx context.Context, in ProductCreateRequest) (*Product, error)
}

type productService struct {
	repo *ProductRepository
}

func NewProductService(repo *ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(ctx context.Context, in ProductCreateRequest) (*Product, error) {
	if in.Name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrValidation)
	}
	if in.Price <= 0 {
		return nil, fmt.Errorf("%w: price must be > 0", ErrValidation)
	}

	p := &Product{
		Name:        in.Name,
		Price:       in.Price,
		Ingredients: pq.StringArray(in.Ingredients),
		Image:       in.Image,
		Rating:      in.Rating,
	}

	created, err := s.repo.Create(p)
	if err != nil {
		return nil, err
	}
	return created, nil
}
