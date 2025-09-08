package mocks

import (
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func NewMockProductDataAccessor() *ProductRepositoryMock {
	return &ProductRepositoryMock{}
}

func (m *ProductRepositoryMock) GetProducts(offset int, limit int, filters models.ProductFilterOptions) ([]models.Product, int, error) {
	args := m.Called(offset, limit, filters)
	return args.Get(0).([]models.Product), args.Int(1), args.Error(2)
}

func (m *ProductRepositoryMock) GetProductDetails(code string) (models.Product, error) {
	args := m.Called(code)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetAllCategories() ([]models.ProductCategory, error) {
	args := m.Called()
	return args.Get(0).([]models.ProductCategory), args.Error(1)
}

func (m *ProductRepositoryMock) CreateCategory(category models.ProductCategory) error {
	args := m.Called(category)
	return args.Error(0)
}
