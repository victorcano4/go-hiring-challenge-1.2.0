package models

import (
	"gorm.io/gorm"
)

type ProductFetcher interface {
	GetPaginatedProducts(int, int) ([]Product, int, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) ProductFetcher {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetPaginatedProducts(offset, limit int) ([]Product, int, error) {
	var products []Product
	var total int64

	if err := r.db.Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Category").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, int(total), nil
}
