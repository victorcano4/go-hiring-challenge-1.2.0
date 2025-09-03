package mocks

import (
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{}
}

func (m *MockProductRepository) GetProducts(offset int, limit int, filters models.ProductFilterOptions) ([]models.Product, int, error) {
	args := m.Called(offset, limit, filters)
	return args.Get(0).([]models.Product), args.Int(1), args.Error(2)
}

func (m *MockProductRepository) GetProductDetails(code string) (models.Product, error) {
	args := m.Called(code)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductRepository) GetAllCategories() ([]models.ProductCategory, error) {
	args := m.Called()
	return args.Get(0).([]models.ProductCategory), args.Error(1)
}
