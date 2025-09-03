package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Products []Product `json:"products"`
}

type Product struct {
	Code     string           `json:"code"`
	Price    float64          `json:"price"`
	Category *ProductCategory `json:"category"`
}

type ProductCategory struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type CatalogHandler struct {
	repo models.ProductFetcher
}

func NewCatalogHandler(r models.ProductFetcher) *CatalogHandler {
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
			http.Error(w, "Invalid priceLessThan parameter", http.StatusBadRequest)
		}
		if parsedPrice > 0 {
			filterOptions.MaxPrice = &parsedPrice
		}
	}

	// Fetch products with filtering and pagination
	res, total, err := h.repo.GetProducts(offset, limit, filterOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Total    int       `json:"total"`
		Products []Product `json:"products"`
	}{
		Total:    total,
		Products: products,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
