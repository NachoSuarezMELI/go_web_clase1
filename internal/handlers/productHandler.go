package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"web/clase1/internal"
	"web/clase1/platform/tools"

	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Service product.ProductService
}

func NewProductHandler(service product.ProductService) *Handler {
	return &Handler{
		Service: service,
	}
}

// GetAllProducts returns all the products in the storage
func (h *Handler) GetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := h.Service.GetAllProducts()
		if err != nil {
			response.JSON(w, http.StatusNotFound, map[string]any{"message": "Products not found"})
		}
		body := product.ResponseProducts{
			Message: "Products found",
			Data:    products,
		}
		response.JSON(w, http.StatusOK, body)
	}
}

// GetProductById returns a product by id
func (h *Handler) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, map[string]any{"message": err.Error()})
			return
		}

		p, err := h.Service.GetProductById(idInt)
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
func (h *Handler) GetProductsByPriceGt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		price := r.URL.Query().Get("priceGt")
		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, map[string]any{"message": err.Error()})
			return
		}

		products := h.Service.FindProductsByPriceGt(priceFloat)
		body := product.ResponseProducts{
			Message: "Products found",
			Data:    products,
		}
		response.JSON(w, http.StatusOK, body)
	}
}

// CreateProduct creates a new product
func (h *Handler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for the token
		if r.Header.Get("Authorization") != os.Getenv("TOKEN") {
			response.JSON(w, http.StatusUnauthorized, map[string]any{"message": "Unauthorized"})
			return
		}

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

		if err = h.Service.CreateProduct(&body); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{"message": err.Error()})
			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{"message": "Product created"})
	}
}

// UpdateOrCreateProduct updates a product or creates it if it doesn't exist
func (h *Handler) UpdateOrCreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for the token
		if r.Header.Get("Authorization") != os.Getenv("TOKEN") {
			response.JSON(w, http.StatusUnauthorized, map[string]any{"message": "Unauthorized"})
			return
		}

		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, map[string]any{"message": err.Error()})
			return
		}

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

		var body product.RequestBodyProduct
		if err := json.Unmarshal(bytes, &body); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{"message": err.Error()})
			return
		}

		if err = h.Service.UpdateOrCreateProduct(&body, idInt); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{"message": err.Error()})
			return
		}

		response.JSON(w, http.StatusNoContent, map[string]any{"message": "Product updated"})
	}
}

// UpdatePartial updates a product partially
func (h *Handler) UpdatePartial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for the token
		if r.Header.Get("Authorization") != os.Getenv("TOKEN") {
			response.JSON(w, http.StatusUnauthorized, map[string]any{"message": "Unauthorized"})
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		bodyMap := make(map[string]any)
		if err := request.JSON(r, &bodyMap); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := h.Service.UpdatePartial(bodyMap, id); err != nil {
			switch {
			case errors.Is(err, product.ErrProdNotFound):
				response.Text(w, http.StatusNotFound, "product not found")
			case errors.Is(err, product.ErrProdInvalidField):
				response.Text(w, http.StatusBadRequest, "invalid field")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")

			}
		}
		response.Text(w, http.StatusNoContent, "product updated")
	}
}

// DeleteProduct deletes a product
func (h *Handler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for the token
		if r.Header.Get("Authorization") != os.Getenv("TOKEN") {
			response.JSON(w, http.StatusUnauthorized, map[string]any{"message": "Unauthorized"})
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		if err := h.Service.DeleteProduct(id); err != nil {
			switch {
			case errors.Is(err, product.ErrProdNotFound):
				response.Text(w, http.StatusNotFound, "product not found")
			case errors.Is(err, product.ErrProdInvalidField):
				response.Text(w, http.StatusBadRequest, "invalid field")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
