package main

import (
	"net/http"
	"web/clase1/internal/product"
	"web/clase1/internal/utils"

	"github.com/go-chi/chi/v5"
)

func main() {

	bytes, err := utils.ReadFile("../docs/db/products.json")
	if err != nil {
		println(err.Error())
	}
	products, err := product.JsonHandler(bytes)
	if err != nil {
		println(err.Error())
	}

	ph := product.NewProductHandler(products)

	r := chi.NewRouter()

	r.Get("/products", ph.GetAllProducts())
	r.Get("/products/{id}", ph.GetProductById())
	r.Get("/products/search", ph.GetProductsByPriceGt())
	r.Post("/products", ph.CreateProduct())

	http.ListenAndServe(":8080", r)
}
