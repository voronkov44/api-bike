package products

import (
	"bike/pkg/db"
	"context"
	"errors"
	"gorm.io/gorm"
)

type ProductRepository struct {
	Database *db.Db
}

func NewProductRepository(database *db.Db) *ProductRepository {
	return &ProductRepository{
		Database: database,
	}
}

func (r *ProductRepository) Create(ctx context.Context, p *Product) (*Product, error) {
	if err := r.Database.DB.WithContext(ctx).Create(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProductRepository) ExistsName(ctx context.Context, name string) (bool, error) {
	var cnt int64
	err := r.Database.DB.WithContext(ctx).Model(&Product{}).
		Where("name = ?", name).Count(&cnt).Error
	return cnt > 0, err
}

func (r *ProductRepository) ExistsSlug(ctx context.Context, slug string) (bool, error) {
	var cnt int64
	err := r.Database.DB.WithContext(ctx).Model(&Product{}).
		Where("slug = ?", slug).Count(&cnt).Error
	return cnt > 0, err
}

func (r *ProductRepository) FindBySlug(ctx context.Context, slug string) (*Product, error) {
	var p Product
	res := r.Database.DB.WithContext(ctx).Where("slug = ?", slug).First(&p)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	return &p, res.Error
}

func (r *ProductRepository) List(ctx context.Context, limit, offset int) ([]Product, error) {
	var list []Product
	q := r.Database.DB.WithContext(ctx).Model(&Product{}).Order("id DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ProductRepository) Save(ctx context.Context, p *Product) (*Product, error) {
	if err := r.Database.DB.WithContext(ctx).Save(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProductRepository) DeleteBySlug(ctx context.Context, slug string) error {
	res := r.Database.DB.WithContext(ctx).Where("slug = ?", slug).Delete(&Product{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
