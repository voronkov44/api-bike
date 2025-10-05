package products

import (
	"bike/pkg/slug"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
)

type ProductService interface {
	Create(ctx context.Context, in ProductCreateRequest) (*Product, error)
	GoTo(ctx context.Context, slug string) (*Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]Product, error)
	Update(ctx context.Context, slug string, in ProductUpdateRequest) (*Product, error)
	ChangeSlug(ctx context.Context, currentSlug, newSlug string) (*Product, error)
	Delete(ctx context.Context, slug string) error
}

type productService struct{ repo *ProductRepository }

func NewProductService(repo *ProductRepository) ProductService { return &productService{repo: repo} }

func (s *productService) Create(ctx context.Context, in ProductCreateRequest) (*Product, error) {
	if in.Name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrValidation)
	}
	if in.Price <= 0 {
		return nil, fmt.Errorf("%w: price must be > 0", ErrValidation)
	}

	// Имя должно быть уникальным
	if ok, err := s.repo.ExistsName(ctx, in.Name); err != nil {
		return nil, err
	} else if ok {
		return nil, fmt.Errorf("%w: name must be unique", ErrValidation)
	}

	// Базовый slug
	base := slug.Slugify(in.Name)
	if base == "" {
		return nil, fmt.Errorf("%w: invalid slug generated from name", ErrValidation)
	}

	// Если slug занят — добавляем короткий uuid-суффикс до уникальности
	use := base
	for {
		exists, err := s.repo.ExistsSlug(ctx, use)
		if err != nil {
			return nil, err
		}
		if !exists {
			break
		}
		use = base + "-" + uuid.NewString()[:8]
	}

	p := &Product{
		Slug:        use,
		Name:        in.Name,
		Type:        in.Type,
		Tags:        pq.StringArray(in.Tags),
		Price:       in.Price,
		Ingredients: pq.StringArray(in.Ingredients),
		Image:       in.Image,
		Rating:      in.Rating,
	}
	return s.repo.Create(ctx, p)
}

func (s *productService) GoTo(ctx context.Context, sl string) (*Product, error) {
	p, err := s.repo.FindBySlug(ctx, sl)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return p, err
}

func (s *productService) GetAll(ctx context.Context, limit, offset int) ([]Product, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *productService) Update(ctx context.Context, sl string, in ProductUpdateRequest) (*Product, error) {
	if in.Name == nil && in.Type == nil && in.Tags == nil &&
		in.Price == nil && in.Ingredients == nil && in.Image == nil && in.Rating == nil {
		return nil, fmt.Errorf("%w: at least one field required", ErrValidation)
	}

	p, err := s.repo.FindBySlug(ctx, sl)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Менять slug автоматически не даём. Имя можно менять — оно уникальное.
	if in.Name != nil && *in.Name != p.Name {
		if ok, err := s.repo.ExistsName(ctx, *in.Name); err != nil {
			return nil, err
		} else if ok {
			return nil, fmt.Errorf("%w: name must be unique", ErrValidation)
		}
		p.Name = *in.Name
	}
	if in.Type != nil {
		p.Type = *in.Type
	}
	if in.Tags != nil {
		p.Tags = pq.StringArray(*in.Tags)
	}
	if in.Price != nil {
		p.Price = *in.Price
	}
	if in.Ingredients != nil {
		p.Ingredients = pq.StringArray(*in.Ingredients)
	}
	if in.Image != nil {
		p.Image = *in.Image
	}
	if in.Rating != nil {
		p.Rating = *in.Rating
	}

	return s.repo.Save(ctx, p)
}

// Явная смена slug пользователем

func (s *productService) ChangeSlug(ctx context.Context, currentSlug, newSlug string) (*Product, error) {
	p, err := s.repo.FindBySlug(ctx, currentSlug)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	ns := slug.Slugify(newSlug)
	if ns == "" {
		return nil, fmt.Errorf("%w: invalid slug", ErrValidation)
	}
	// Проверяем, что новый slug свободен
	if ok, err := s.repo.ExistsSlug(ctx, ns); err != nil {
		return nil, err
	} else if ok {
		return nil, fmt.Errorf("%w: slug already exists", ErrValidation)
	}

	p.Slug = ns
	return s.repo.Save(ctx, p)
}

func (s *productService) Delete(ctx context.Context, sl string) error {
	err := s.repo.DeleteBySlug(ctx, sl)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
