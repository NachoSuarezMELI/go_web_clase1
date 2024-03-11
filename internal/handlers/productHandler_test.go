package handlers

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"web/clase1/internal/repository"
	"web/clase1/internal/service"
	"web/clase1/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestGetProduct(t *testing.T) {
	t.Run("should return all products", func(t *testing.T) {
		// Arrange
		st := storage.NewStorageJSON("../../docs/db/test_products.json")
		rp := repository.NewProductRepository(st)
		sv := service.NewProductService(rp)
		hd := NewProductHandler(sv)
		// Act
		req := httptest.NewRequest("GET", "/products", nil)
		res := httptest.NewRecorder()
		hdFunc := hd.GetAllProducts()
		hdFunc(res, req)
		// Assert
		expectedBody := `{"status_code":200,"message":"Products found","data":[{"id":1,"name":"Oil - Margarine","quantity":439,"code_value":"S82254D","is_published":true,"expiration":"15/12/2021","price":71.42},{"id":2,"name":"Pineapple - Canned, Rings","quantity":345,"code_value":"M4637","is_published":true,"expiration":"09/08/2021","price":352.79},{"id":3,"name":"Wine - Red Oakridge Merlot","quantity":367,"code_value":"T65812","is_published":false,"expiration":"24/05/2021","price":179.23}]}`

		require.Equal(t, 200, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
}

func TestGetProductById(t *testing.T) {
	t.Run("should return a product by id", func(t *testing.T) {
		// Arrange
		st := storage.NewStorageJSON("../../docs/db/test_products.json")
		rp := repository.NewProductRepository(st)
		sv := service.NewProductService(rp)
		hd := NewProductHandler(sv)
		// Act
		req := httptest.NewRequest("GET", "/products/1", nil)
		// Set query params with context
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		res := httptest.NewRecorder()
		hdFunc := hd.GetProductById()
		hdFunc(res, req)
		// Assert
		expectedBody := `{"status_code":200,"message":"Product found","data":{"id":1,"name":"Oil - Margarine","quantity":439,"code_value":"S82254D","is_published":true,"expiration":"15/12/2021","price":71.42}}`

		require.Equal(t, 200, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
	t.Run("should return a bad request when id is not a number", func(t *testing.T) {
		// Arrange
		st := storage.NewStorageJSON("../../docs/db/test_products.json")
		rp := repository.NewProductRepository(st)
		sv := service.NewProductService(rp)
		hd := NewProductHandler(sv)
		// Act
		req := httptest.NewRequest("GET", "/products/abc", nil)
		// Set query params with context
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "abc")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		res := httptest.NewRecorder()
		hdFunc := hd.GetProductById()
		hdFunc(res, req)
		// Assert
		expectedBody := `{"status_code":400,"message":"Bad request","data":null}`

		require.Equal(t, 400, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
	t.Run("should return a not found when id is not found", func(t *testing.T) {
		// Arrange
		st := storage.NewStorageJSON("../../docs/db/test_products.json")
		rp := repository.NewProductRepository(st)
		sv := service.NewProductService(rp)
		hd := NewProductHandler(sv)
		// Act
		req := httptest.NewRequest("GET", "/products/4", nil)
		// Set query params with context
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "4")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		res := httptest.NewRecorder()
		hdFunc := hd.GetProductById()
		hdFunc(res, req)
		// Assert
		expectedBody := `{"status_code":404,"message":"product not found","data":null}`

		require.Equal(t, 404, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
}

func TestCreateProduct(t *testing.T) {
	t.Run("should create a product", func(t *testing.T) {
		// Arrange
		st := storage.NewStorageJSON("../../docs/db/test_products.json")
		rp := repository.NewProductRepository(st)
		sv := service.NewProductService(rp)
		hd := NewProductHandler(sv)

		body := `{"name":"Oil - Margarine","quantity":439,"code_value":"S82254D","is_published":true,"expiration":"15/12/2021","price":71.42}`

		// Act
		req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
		res := httptest.NewRecorder()
		hdFunc := hd.CreateProduct()
		hdFunc(res, req)

		// Assert
		expectedBody := `{"status_code":201,"message":"Product created","data":{"id":1,"name":"Oil - Margarine","quantity":439,"code_value":"S82254D","is_published":true,"expiration":"15/12/2021","price":71.42}}`

		require.Equal(t, 201, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("should delete a product", func(t *testing.T) {
		// Arrange
		st := storage.NewStorageJSON("../../docs/db/test_products.json")
		rp := repository.NewProductRepository(st)
		sv := service.NewProductService(rp)
		hd := NewProductHandler(sv)

		// Act
		req := httptest.NewRequest("DELETE", "/products/2", nil)
		// Set query params with context
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		res := httptest.NewRecorder()
		hdFunc := hd.DeleteProduct()
		hdFunc(res, req)

		// Assert
		expectedBody := `{"status_code":204,"message":"product deleted","data":null}`

		require.Equal(t, 204, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
	t.Run("should return a bad request when id is not a number", func(t *testing.T) {
		// Arrange
		st := storage.NewStorageJSON("../../docs/db/test_products.json")
		rp := repository.NewProductRepository(st)
		sv := service.NewProductService(rp)
		hd := NewProductHandler(sv)
		// Act
		req := httptest.NewRequest("DELETE", "/products/abc", nil)
		// Set query params with context
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "abc")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		res := httptest.NewRecorder()
		hdFunc := hd.DeleteProduct()
		hdFunc(res, req)
		// Assert
		expectedBody := `{"status_code":400,"message":"invalid id","data":null}`

		require.Equal(t, 400, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
	t.Run("should return a not found when id is not found", func(t *testing.T) {
		// Arrange
		st := storage.NewStorageJSON("../../docs/db/test_products.json")
		rp := repository.NewProductRepository(st)
		sv := service.NewProductService(rp)
		hd := NewProductHandler(sv)
		// Act
		req := httptest.NewRequest("DELETE", "/products/4", nil)
		// Set query params with context
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "4")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		res := httptest.NewRecorder()
		hdFunc := hd.DeleteProduct()
		hdFunc(res, req)
		// Assert
		expectedBody := `{"status_code":404,"message":"product not found","data":null}`

		require.Equal(t, 404, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
}

func TestUpdateOrCreateProduct(t *testing.T) {
	t.Run("should throw a bad request when id is not a number", func(t *testing.T) {
		// Arrange
		st := storage.NewStorageJSON("../../docs/db/test_products.json")
		rp := repository.NewProductRepository(st)
		sv := service.NewProductService(rp)
		hd := NewProductHandler(sv)
		// Act
		req := httptest.NewRequest("PUT", "/products/abc", nil)
		// Set query params with context
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "abc")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		res := httptest.NewRecorder()
		hdFunc := hd.UpdateOrCreateProduct()
		hdFunc(res, req)
		// Assert
		expectedBody := `{"status_code":400,"message":"invalid id","data":null}`

		require.Equal(t, 400, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
}

func TestUpdatePartial(t *testing.T) {
	t.Run("should throw a bad request when id is not a number", func(t *testing.T) {
		// Arrange
		st := storage.NewStorageJSON("../../docs/db/test_products.json")
		rp := repository.NewProductRepository(st)
		sv := service.NewProductService(rp)
		hd := NewProductHandler(sv)
		// Act
		req := httptest.NewRequest("PATCH", "/products/abc", nil)
		// Set query params with context
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "abc")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		req.Header.Set("Content-Type", "application/json")

		res := httptest.NewRecorder()
		hdFunc := hd.UpdatePartial()
		hdFunc(res, req)
		// Assert
		expectedBody := `{"status_code":400,"message":"invalid id","data":null}`
		require.Equal(t, 400, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
	t.Run("should throw a not found when id is not found", func(t *testing.T) {
		// Arrange
		st := storage.NewStorageJSON("../../docs/db/test_products.json")
		rp := repository.NewProductRepository(st)
		sv := service.NewProductService(rp)
		hd := NewProductHandler(sv)

		body := `{"price": 99999}`

		// Act
		req := httptest.NewRequest("PATCH", "/products/4", strings.NewReader(body))
		// Set query params with context
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "4")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		res := httptest.NewRecorder()
		hdFunc := hd.UpdatePartial()
		hdFunc(res, req)
		// Assert
		expectedBody := `{"status_code":404,"message":"product not found","data":null}`
		require.Equal(t, 404, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
}
