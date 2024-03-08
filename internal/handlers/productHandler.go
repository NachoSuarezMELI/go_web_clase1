package handlers

import (
	"encoding/json"
	"errors"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"web/clase1/internal"
	"web/clase1/internal/service"
	"web/clase1/platform/tools"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// GetAllProducts returns all the products in the storage
func (h *ProductHandler) GetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products := h.service.GetAllProducts()
		body := product.ResponseProducts{
			Message: "Products found",
			Data:    products,
		}
		response.JSON(w, http.StatusOK, body)
	}
}

// GetProductById returns a product by id
func (h *ProductHandler) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, map[string]any{"message": err.Error()})
			return
		}

		p, err := h.service.GetProductById(idInt)
		if err != nil {
			response.JSON(w, http.StatusNotFound, map[string]any{"message": "Product not found"})
			return
		}
		body := product.ResponseProduct{
			Message: "Product found",
			Data:    *p,
		}
		response.JSON(w, http.StatusOK, body)
	}
}

// GetProductsByPriceGt returns a list of products with a price greater than the one specified in the query
func (h *ProductHandler) GetProductsByPriceGt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		price := r.URL.Query().Get("priceGt")
		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, map[string]any{"message": err.Error()})
			return
		}

		products := h.service.FindProductsByPriceGt(priceFloat)
		body := product.ResponseProducts{
			Message: "Products found",
			Data:    products,
		}
		response.JSON(w, http.StatusOK, body)
	}
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{"message": err.Error()})
			return
		}

		bodyMap := make(map[string]any)
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{"message": err.Error()})
			return
		}

		if err := tools.CheckFieldExistance(bodyMap, "name", "quantity", "code_value", "is_published", "expiration", "price"); err != nil {
			var fieldError *tools.FieldError
			if errors.As(err, &fieldError) {
				response.JSON(w, http.StatusBadRequest, map[string]any{"message": err.Error()})
				return
			}

			response.JSON(w, http.StatusBadRequest, map[string]any{"message": err.Error()})
			return
		}

		var body product.Product
		if err := json.Unmarshal(bytes, &body); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{"message": err.Error()})
			return
		}

		if err = h.service.CreateProduct(&body); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{"message": err.Error()})
			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{"message": "Product created"})
	}
}
