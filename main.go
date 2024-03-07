package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"web/clase1/handlers"
	"web/clase1/services"
)

func main() {

	bytes, err := services.ReadFile("products.json")
	if err != nil {
		println(err.Error())
	}
	products, err := services.JsonHandler(bytes)
	if err != nil {
		println(err.Error())
	}

	ph := controllers.NewProductHandler(products)

	r := chi.NewRouter()

	r.Get("/products", ph.GetAllProducts())
	r.Get("/products/{id}", ph.GetProductById())
	r.Get("/products/search", ph.GetProductsByPriceGt())
	r.Post("/products", ph.CreateProduct())

	http.ListenAndServe(":8080", r)
}
