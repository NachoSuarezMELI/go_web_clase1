package product

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"web/clase1/internal/web"
)

type ProductHandler struct {
	storage map[int]*Product
}

func NewProductHandler(products []Product) *ProductHandler {
	storage := make(map[int]*Product)
	for _, product := range products {
		storage[product.Id] = &product
	}
	return &ProductHandler{
		storage: storage,
	}
}

// GetAllProducts returns all the products in the storage
func (h *ProductHandler) GetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(h.storage) == 0 {
			web.ResponseJson(w, map[string]any{"message": "No products found"}, http.StatusNotFound)
			return
		}
		products := GetProducts(h.storage)

		body := ResponseProducts{
			Message: "Products found",
			Data:    products,
		}

		web.ResponseJson(w, body, http.StatusOK)
	}
}

// GetProductById returns a product by id
func (h *ProductHandler) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			web.ResponseJson(w, map[string]any{"message": err.Error()}, http.StatusInternalServerError)
			return
		}
		p := FindProductById(idInt, h.storage)
		if p == nil {
			web.ResponseJson(w, map[string]any{"message": "Product not found"}, http.StatusNotFound)
			return
		}

		body := ResponseProduct{
			Message: "Product found",
			Data:    *p,
		}
		web.ResponseJson(w, body, http.StatusOK)
	}
}

// GetProductsByPriceGt returns a list of products with a price greater than the one specified in the query
func (h *ProductHandler) GetProductsByPriceGt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		priceGt := r.URL.Query().Get("priceGt")
		priceGtFloat, err := strconv.ParseFloat(priceGt, 64)

		if err != nil {
			web.ResponseJson(w, map[string]any{"message": err.Error()}, http.StatusInternalServerError)
			return
		}

		products := FindProductsByPriceGt(priceGtFloat, h.storage)
		if len(products) == 0 {
			web.ResponseJson(w, map[string]any{"message": "No products found"}, http.StatusNotFound)
			return
		}

		body := ResponseProducts{
			Message: "Products found",
			Data:    products,
		}

		web.ResponseJson(w, body, http.StatusOK)
	}
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Requestbody RequestBodyProduct
		err := web.RequestJsonProduct(r, &Requestbody)
		if err != nil {
			web.ResponseJson(w, map[string]any{"message": err.Error()}, http.StatusBadRequest)
			return
		}

		p, err := CreateProduct(len(h.storage)+1, Requestbody, h.storage)
		if err != nil {
			web.ResponseJson(w, map[string]any{"message": err.Error()}, http.StatusBadRequest)
			return
		}
		Responsebody := ResponseProduct{
			Message: "Product created",
			Data:    *p,
		}

		web.ResponseJson(w, Responsebody, http.StatusCreated)
	}
}
