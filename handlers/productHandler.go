package controllers

import (
	"net/http"
	"strconv"
	"web/clase1/models"
	"web/clase1/services"
	"web/clase1/web"

	"github.com/go-chi/chi/v5"
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
			web.ResponseJson(w, map[string]any{"message": "No products found"}, http.StatusNotFound)
			return
		}

		products := services.GetProducts(h.storage)

		body := models.ResponseProducts{
			Message: "Products found",
			Data:    products,
		}

		web.ResponseJson(w, body, http.StatusOK)
	}
}

func (h *ProductHandler) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			web.ResponseJson(w, map[string]any{"message": err.Error()}, http.StatusInternalServerError)
			return
		}
		product := services.FindProductById(idInt, h.storage)
		if product == nil {
			web.ResponseJson(w, map[string]any{"message": "Product not found"}, http.StatusNotFound)
			return
		}

		body := models.ResponseProduct{
			Message: "Product found",
			Data:    *product,
		}
		web.ResponseJson(w, body, http.StatusOK)
	}
}

func (h *ProductHandler) GetProductsByPriceGt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		priceGt := r.URL.Query().Get("priceGt")
		priceGtFloat, err := strconv.ParseFloat(priceGt, 64)

		if err != nil {
			web.ResponseJson(w, map[string]any{"message": err.Error()}, http.StatusInternalServerError)
			return
		}

		products := services.FindProductsByPriceGt(priceGtFloat, h.storage)
		if len(products) == 0 {
			web.ResponseJson(w, map[string]any{"message": "No products found"}, http.StatusNotFound)
			return
		}

		body := models.ResponseProducts{
			Message: "Products found",
			Data:    products,
		}

		web.ResponseJson(w, body, http.StatusOK)
	}
}

func (h *ProductHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Requestbody models.RequestBodyProduct
		err := web.RequestJsonProduct(r, &Requestbody)
		if err != nil {
			web.ResponseJson(w, map[string]any{"message": err.Error()}, http.StatusBadRequest)
			return
		}

		p, err := services.CreateProduct(len(h.storage)+1, Requestbody, h.storage)
		if err != nil {
			web.ResponseJson(w, map[string]any{"message": err.Error()}, http.StatusBadRequest)
			return
		}
		Responsebody := models.ResponseProduct{
			Message: "Product created",
			Data:    *p,
		}

		web.ResponseJson(w, Responsebody, http.StatusCreated)
	}
}
