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
		if err != nil {
			http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
		}
		if parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		}
		if parsedLimit > 100 {
			limit = 100
		} else if parsedLimit >= 1 {
			limit = parsedLimit
		}
	}

	res, total, err := h.repo.GetPaginatedProducts(offset, limit)
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
