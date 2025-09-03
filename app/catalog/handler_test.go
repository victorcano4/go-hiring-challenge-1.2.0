package catalog_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	catalog "github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/mocks"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
	// Import the mock package
)

func TestHandleGetCategories(t *testing.T) {
	mockRepo := mocks.NewMockProductDataAccessor()
	handler := catalog.NewCatalogHandler(mockRepo)

	t.Run("Categories found", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil

		categories := []models.ProductCategory{
			{Code: "CLOTHING", Name: "Clothing"},
			{Code: "SHOES", Name: "Shoes"},
			{Code: "ACCESSORIES", Name: "Accessories"},
		}

		mockRepo.On("GetAllCategories").Return(categories, nil).Once()

		r := mux.NewRouter()
		r.HandleFunc("/categories", handler.HandleGetCategories)

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response []catalog.ProductCategory
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 3)
		assert.Equal(t, "CLOTHING", response[0].Code)
		assert.Equal(t, "Clothing", response[0].Name)
		assert.Equal(t, "SHOES", response[1].Code)
		assert.Equal(t, "Shoes", response[1].Name)
		assert.Equal(t, "ACCESSORIES", response[2].Code)
		assert.Equal(t, "Accessories", response[2].Name)

		mockRepo.AssertExpectations(t)
	})

	t.Run("No categories found", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil

		mockRepo.On("GetAllCategories").Return([]models.ProductCategory{}, nil).Once()

		r := mux.NewRouter()
		r.HandleFunc("/categories", handler.HandleGetCategories)

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response []catalog.ProductCategory
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error fetching categories", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil

		mockRepo.On("GetAllCategories").Return([]models.ProductCategory{}, errors.New("database error")).Once()

		r := mux.NewRouter()
		r.HandleFunc("/categories", handler.HandleGetCategories)

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, "database error\n", resp.Body.String())
		mockRepo.AssertExpectations(t)
	})
}

func TestHandleCreateCategory(t *testing.T) {
	mockRepo := mocks.NewMockProductDataAccessor()
	handler := catalog.NewCatalogHandler(mockRepo)

	t.Run("Successfully create category", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil

		newCategory := models.ProductCategory{Code: "NEW", Name: "New Category"}
		mockRepo.On("CreateCategory", newCategory).Return(nil).Once()

		r := mux.NewRouter()
		r.HandleFunc("/categories", handler.HandleCreateCategory).Methods(http.MethodPost)

		body, _ := json.Marshal(newCategory)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Equal(t, "{\"code\":\"NEW\",\"name\":\"New Category\"}\n", resp.Body.String())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error creating category", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil

		newCategory := models.ProductCategory{Code: "NEW", Name: "New Category"}
		mockRepo.On("CreateCategory", newCategory).Return(errors.New("database error")).Once()

		r := mux.NewRouter()
		r.HandleFunc("/categories", handler.HandleCreateCategory).Methods(http.MethodPost)

		body, _ := json.Marshal(newCategory)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, "database error\n", resp.Body.String())
		mockRepo.AssertExpectations(t)
	})
}
