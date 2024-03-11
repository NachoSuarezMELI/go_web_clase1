package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	product "web/clase1/internal"
	"web/clase1/internal/web"
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
			body := web.StandarResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
			}
			response.JSON(w, http.StatusNotFound, body)
		}

		body := web.StandarResponse{
			StatusCode: http.StatusOK,
			Message:    "Products found",
			Data:       products,
		}
		response.JSON(w, http.StatusOK, body)
	}
}

// GetProductById returns a product by id
func (h *Handler) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for the token
		if r.Header.Get("Authorization") != os.Getenv("TOKEN") {
			body := web.StandarResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Unauthorized",
			}
			response.JSON(w, http.StatusUnauthorized, body)
			return
		}
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Bad request",
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		p, err := h.Service.GetProductById(idInt)
		if err != nil {
			if errors.Is(err, product.ErrProdNotFound) {
				body := web.StandarResponse{
					StatusCode: http.StatusNotFound,
					Message:    err.Error(),
				}
				response.JSON(w, http.StatusNotFound, body)
				return
			} else {
				body := web.StandarResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    err.Error(),
				}
				response.JSON(w, http.StatusInternalServerError, body)
				return
			}
		}
		body := web.StandarResponse{
			StatusCode: http.StatusOK,
			Message:    "Product found",
			Data:       p,
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
		body := web.StandarResponse{
			StatusCode: http.StatusOK,
			Message:    "Products found",
			Data:       products,
		}
		response.JSON(w, http.StatusOK, body)
	}
}

// CreateProduct creates a new product
func (h *Handler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for the token
		if r.Header.Get("Authorization") != os.Getenv("TOKEN") {
			body := web.StandarResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Unauthorized",
			}
			response.JSON(w, http.StatusUnauthorized, body)
			return
		}

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		bodyMap := make(map[string]any)
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		if err := tools.CheckFieldExistance(bodyMap, "name", "quantity", "code_value", "is_published", "expiration", "price"); err != nil {
			var fieldError *tools.FieldError
			if errors.As(err, &fieldError) {
				body := web.StandarResponse{
					StatusCode: http.StatusBadRequest,
					Message:    err.Error(),
				}
				response.JSON(w, http.StatusBadRequest, body)
				return
			}

			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		var p product.Product
		if err := json.Unmarshal(bytes, &p); err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		if err = h.Service.CreateProduct(&p); err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		body := web.StandarResponse{
			StatusCode: http.StatusCreated,
			Message:    "Product created",
			Data:       p,
		}
		response.JSON(w, http.StatusCreated, body)
	}
}

// UpdateOrCreateProduct updates a product or creates it if it doesn't exist
func (h *Handler) UpdateOrCreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for the token
		if r.Header.Get("Authorization") != os.Getenv("TOKEN") {
			body := web.StandarResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Unauthorized",
			}
			response.JSON(w, http.StatusUnauthorized, body)
			return
		}

		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid id",
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		bodyMap := make(map[string]any)
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		if err := tools.CheckFieldExistance(bodyMap, "name", "quantity", "code_value", "is_published", "expiration", "price"); err != nil {
			var fieldError *tools.FieldError
			if errors.As(err, &fieldError) {
				body := web.StandarResponse{
					StatusCode: http.StatusBadRequest,
					Message:    err.Error(),
				}
				response.JSON(w, http.StatusBadRequest, body)
				return
			}

			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}

			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		var p product.RequestBodyProduct
		if err := json.Unmarshal(bytes, &p); err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		if err = h.Service.UpdateOrCreateProduct(&p, idInt); err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		product, _ := h.Service.GetProductById(idInt)

		body := web.StandarResponse{
			StatusCode: http.StatusNoContent,
			Message:    "Product updated",
			Data:       product,
		}

		response.JSON(w, http.StatusNoContent, body)
	}
}

// UpdatePartial updates a product partially
func (h *Handler) UpdatePartial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for the token
		if r.Header.Get("Authorization") != os.Getenv("TOKEN") {
			body := web.StandarResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Unauthorized",
			}
			response.JSON(w, http.StatusUnauthorized, body)
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid id",
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		bodyMap := make(map[string]any)
		if err := request.JSON(r, &bodyMap); err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid request body",
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		if err := h.Service.UpdatePartial(bodyMap, id); err != nil {
			switch {
			case errors.Is(err, product.ErrProdNotFound):
				body := web.StandarResponse{
					StatusCode: http.StatusNotFound,
					Message:    "product not found",
				}
				response.JSON(w, http.StatusNotFound, body)
			case errors.Is(err, product.ErrProdInvalidField):
				body := web.StandarResponse{
					StatusCode: http.StatusBadRequest,
					Message:    err.Error(),
				}
				response.JSON(w, http.StatusBadRequest, body)
			default:
				body := web.StandarResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "internal server error",
				}
				response.JSON(w, http.StatusInternalServerError, body)
			}
		}
		body := web.StandarResponse{
			StatusCode: http.StatusNoContent,
			Message:    "product updated",
		}
		response.JSON(w, http.StatusNoContent, body)
	}
}

// DeleteProduct deletes a product
func (h *Handler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for the token
		if r.Header.Get("Authorization") != os.Getenv("TOKEN") {
			body := web.StandarResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Unauthorized",
			}
			response.JSON(w, http.StatusUnauthorized, body)
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid id",
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		if err := h.Service.DeleteProduct(id); err != nil {
			switch {
			case errors.Is(err, product.ErrProdNotFound):
				body := web.StandarResponse{
					StatusCode: http.StatusNotFound,
					Message:    "product not found",
				}
				response.JSON(w, http.StatusNotFound, body)
			case errors.Is(err, product.ErrProdInvalidField):
				body := web.StandarResponse{
					StatusCode: http.StatusBadRequest,
					Message:    err.Error(),
				}
				response.JSON(w, http.StatusBadRequest, body)
			default:
				body := web.StandarResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "internal server error",
				}
				response.JSON(w, http.StatusInternalServerError, body)
			}
			return
		}

		body := web.StandarResponse{
			StatusCode: http.StatusNoContent,
			Message:    "product deleted",
		}

		response.JSON(w, http.StatusNoContent, body)
	}
}

// GetConsumerPrice returns the consumer price of a list of products
// TODO: fix update quantity--
func (h *Handler) GetConsumerPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for the token
		if r.Header.Get("Authorization") != os.Getenv("TOKEN") {
			body := web.StandarResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Unauthorized",
			}
			response.JSON(w, http.StatusUnauthorized, body)
			return
		}

		strIds := r.URL.Query()["list"]
		var ids []int

		err := json.Unmarshal([]byte(strIds[0]), &ids)
		if err != nil {
			body := web.StandarResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid id",
			}
			response.JSON(w, http.StatusBadRequest, body)
			return
		}

		var price float64

		if len(ids) == 0 {
			products, err := h.Service.GetAllProducts()
			if err != nil {
				body := web.StandarResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "internal server error",
				}
				response.JSON(w, http.StatusInternalServerError, body)
				return
			} else {
				for _, product := range products {
					if product.Quantity > 0 && product.Is_Published {
						price += product.Price
						product.Quantity--
					} else {
						body := web.StandarResponse{
							StatusCode: http.StatusBadRequest,
							Message:    "product not available",
						}
						response.JSON(w, http.StatusBadRequest, body)
						return
					}
					price = price * 1.15
					body := web.StandarResponse{
						StatusCode: http.StatusOK,
						Message:    "consumer price",
						Data:       fmt.Sprintf("%.2f", price),
					}
					response.JSON(w, http.StatusOK, body)
					return
				}
			}
		}

		// Fix product quantity --
		// Using GetProductById method doesn't update the quantity of the product
		// Fix it using the UpdatePartial method?

		for _, id := range ids {
			product, err := h.Service.GetProductById(id)
			if err != nil {
				body := web.StandarResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "internal server error",
				}
				response.JSON(w, http.StatusInternalServerError, body)
				return
			} else {
				if product.Quantity > 0 && product.Is_Published {
					price += product.Price

					//
					product.Quantity--
					//
				} else {
					body := web.StandarResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "product not available",
					}
					response.JSON(w, http.StatusBadRequest, body)
					return
				}
			}
			price += product.Price
		}
		switch {
		case len(ids) < 10:
			price = price * 1.21
		case len(ids) >= 10 && len(ids) < 20:
			price = price * 1.17
		case len(ids) >= 20:
			price = price * 1.15
		}
		body := web.StandarResponse{
			StatusCode: http.StatusOK,
			Message:    "consumer price",
			Data:       fmt.Sprintf("%.2f", price),
		}
		response.JSON(w, http.StatusOK, body)
	}
}
