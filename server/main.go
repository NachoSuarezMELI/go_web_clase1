package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"web/clase1/internal/handlers"
	"web/clase1/internal/repository"
	"web/clase1/internal/service"
)

func main() {

	rp := repository.NewProductRepository(nil, 0)
	sv := service.NewProductService(rp)
	h := handlers.NewProductHandler(sv)

	if err := rp.GeneretaDB(); err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	router.Get("/products", h.GetAllProducts())
	router.Get("/products/{id}", h.GetProductById())
	router.Get("/products/search", h.GetProductsByPriceGt())
	router.Post("/products", h.CreateProduct())
	router.Put("/products/{id}", h.UpdateOrCreateProduct())

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
