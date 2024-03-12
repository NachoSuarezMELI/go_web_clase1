package main

import (
	"net/http"
	"os"
	"web/clase1/internal/handlers"
	"web/clase1/internal/repository"
	"web/clase1/internal/service"
	"web/clase1/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	os.Setenv("TOKEN", "123456")

	st := storage.NewStorageJSON("../docs/db/products.json")
	rp := repository.NewProductRepository(st)
	sv := service.NewProductService(rp)
	h := handlers.NewProductHandler(sv)

	router := chi.NewRouter()

	router.Get("/products", h.GetAllProducts())
	router.Get("/products/{id}", h.GetProductById())
	router.Get("/products/search", h.GetProductsByPriceGt())
	router.Post("/products", h.CreateProduct())
	router.Put("/products/{id}", h.UpdateOrCreateProduct())
	router.Patch("/products/{id}", h.UpdatePartial())
	router.Delete("/products/{id}", h.DeleteProduct())
	router.Get("/products/consumer_price", h.GetConsumerPrice())

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
