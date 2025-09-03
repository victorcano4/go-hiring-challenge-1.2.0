package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
)

type Response struct {
	Products []Product `json:"products"`
}

type Product struct {
	Code     string           `json:"code"`
	Price    float64          `json:"price"`
	Category *ProductCategory `json:"category"`
	Variants []Variant        `json:"variants"`
}

type ProductCategory struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Variant struct {
	Name  string  `json:"name"`
	Sku   string  `json:"sku"`
	Price float64 `json:"price"`
}

type CatalogHandler struct {
	repo models.ProductDataAccessor
}

func NewCatalogHandler(r models.ProductDataAccessor) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	offsetParam := r.URL.Query().Get("offset")
	limitParam := r.URL.Query().Get("limit")

	offset := 0
	limit := 10

	if offsetParam != "" {
		parsedOffset, err := strconv.Atoi(offsetParam)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err == nil {
			if parsedLimit > 100 || parsedLimit < 1 {
				limit = 100
			} else {
				limit = parsedLimit
			}
		} else {
			parsedLimit = 100
		}
	}

	// Parse query parameters for filtering
	categoryParam := r.URL.Query().Get("category")
	priceLessThanParam := r.URL.Query().Get("priceLessThan")

	// Initialize filtering options
	filterOptions := models.ProductFilterOptions{}

	if categoryParam != "" {
		filterOptions.Category = &categoryParam
	}

	if priceLessThanParam != "" {
		parsedPrice, err := strconv.ParseFloat(priceLessThanParam, 64)
		if err != nil {
			api.ErrorResponse(w, http.StatusBadRequest, "Invalid priceLessThan parameter")
			return
		}
		if parsedPrice > 0 {
			filterOptions.MaxPrice = &parsedPrice
		}
	}

	// Fetch products with filtering and pagination
	res, total, err := h.repo.GetProducts(offset, limit, filterOptions)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map response
	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
		}
		if p.Category.ID != 0 {
			products[i].Category = &ProductCategory{
				Name: p.Category.Name,
				Code: p.Category.Code,
			}
		}
	}

	// Return the products as a JSON response
	response := struct {
		Total    int       `json:"total"`
		Products []Product `json:"products"`
	}{
		Total:    total,
		Products: products,
	}

	api.OKResponse(w, response)
}

func (h *CatalogHandler) HandleGetProductDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	codeParam := vars["code"]

	product, err := h.repo.GetProductDetails(codeParam)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.ErrorResponse(w, http.StatusNotFound, "Product not found")
			return
		}
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := Product{
		Code:     product.Code,
		Price:    product.Price.InexactFloat64(),
		Category: &ProductCategory{Name: product.Category.Name, Code: product.Category.Code},
		Variants: make([]Variant, len(product.Variants)),
	}

	for i, v := range product.Variants {
		variantPrice := v.Price
		if variantPrice.IsZero() {
			variantPrice = product.Price
		}
		response.Variants[i] = Variant{
			Name:  v.Name,
			Sku:   v.SKU,
			Price: variantPrice.InexactFloat64(),
		}
	}

	api.OKResponse(w, response)
}

// HandleGetCategories handles the categories endpoint
func (h *CatalogHandler) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.GetAllCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map categories to response
	response := make([]ProductCategory, len(categories))
	for i, c := range categories {
		response[i] = ProductCategory{
			Code: c.Code,
			Name: c.Name,
		}
	}

	api.OKResponse(w, response)
}

func (h *CatalogHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory models.ProductCategory

	// Parse the JSON body
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate the category details
	if newCategory.Code == "" || newCategory.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Category code and name are required")
		return
	}

	// Insert the new category into the database
	if err := h.repo.CreateCategory(newCategory); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.CreatedResponse(w, newCategory)
}
