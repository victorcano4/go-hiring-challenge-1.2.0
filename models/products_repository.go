package models

import (
	"gorm.io/gorm"
)

type ProductFetcher interface {
	GetAllProducts() ([]Product, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) ProductFetcher {
	return ProductsRepository{
		db: db,
	}
}

func (pr ProductsRepository) GetAllProducts() ([]Product, error) {
	var products []Product
	if err := pr.db.Preload("Variants").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
