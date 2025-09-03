package models

import (
	"gorm.io/gorm"
)

type ProductFetcher interface {
	GetProducts(int, int, ProductFilterOptions) ([]Product, int, error)
	GetProductDetails(string) (Product, error)
	GetAllCategories() ([]ProductCategory, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) ProductFetcher {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetProducts(offset, limit int, filters ProductFilterOptions) ([]Product, int, error) {
	var products []Product
	var total int64

	if err := r.db.Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	paginatedProducts := r.db.Preload("Category").Offset(offset).Limit(limit)
	if filters.MaxPrice != nil {
		paginatedProducts = paginatedProducts.Where("price < ?", *filters.MaxPrice)
	}
	if filters.Category != nil {
		paginatedProducts = paginatedProducts.Joins("JOIN product_categories ON product_categories.id = products.category_id").Where("product_categories.code = ?", *filters.Category)
	}

	if err := paginatedProducts.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, int(total), nil
}

// GetProductDetails fetches product details by code, including its variants and category
func (r *ProductsRepository) GetProductDetails(code string) (Product, error) {
	var product Product
	if err := r.db.Preload("Variants").Preload("Category").Where("code = ?", code).First(&product).Error; err != nil {
		return product, err
	}
	return product, nil
}

// GetAllCategories fetches all product categories from the database
func (r *ProductsRepository) GetAllCategories() ([]ProductCategory, error) {
	var categories []ProductCategory
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
