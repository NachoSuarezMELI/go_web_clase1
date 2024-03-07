package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"web/clase1/models"
	"web/clase1/services"
)

type ProductHandler struct {
	storage map[int]*models.Product
}

func NewProductHandler(products []models.Product) *ProductHandler {
	storage := make(map[int]*models.Product)
	for _, product := range products {
		storage[product.Id] = &product
	}
	return &ProductHandler{
		storage: storage,
	}
}

func (h *ProductHandler) GetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(h.storage) == 0 {
			http.Error(w, "No products found", http.StatusNotFound)
			return
		}

		products := services.GetProducts(h.storage)

		body := models.ResponseProducts{
			Message: "Products found",
			Data:    products,
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}
}

func (h *ProductHandler) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		product := services.FindProductById(idInt, h.storage)
		if product == nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		body := models.ResponseProduct{
			Message: "Product found",
			Data:    *product,
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}
}

func (h *ProductHandler) GetProductsByPriceGt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		priceGt := r.URL.Query().Get("priceGt")
		priceGtFloat, err := strconv.ParseFloat(priceGt, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		products := services.FindProductsByPriceGt(priceGtFloat, h.storage)
		if len(products) == 0 {
			http.Error(w, "No products found", http.StatusNotFound)
			return
		}

		body := models.ResponseProducts{
			Message: "Products found",
			Data:    products,
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}
}

func (h *ProductHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body models.RequestBodyProduct
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p, err := services.CreateProduct(len(h.storage)+1, body, h.storage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		Responsebody := models.ResponseProduct{
			Message: "Product created",
			Data:    *p,
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Responsebody)
	}
}
