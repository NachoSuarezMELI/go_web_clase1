package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"web/clase1/models"
	"web/clase1/services"

	"github.com/go-chi/chi/v5"
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

	r := chi.NewRouter()

	r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		w.Write(bytes)
	})

	r.Get("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if idInt <= 0 {
			http.Error(w, "The product id cannot be zero or lower", http.StatusBadRequest)
			return
		}

		var foundProduct *models.Product
		foundProduct = services.FindProductById(idInt, products)

		if foundProduct == nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		response, err := json.Marshal(foundProduct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

	//Crear una ruta /products/search que nos permita buscar por parámetro los productos cuyo precio sean mayor a un valor priceGt.
	r.Get("/products/search", func(w http.ResponseWriter, r *http.Request) {
		priceGt := r.URL.Query().Get("priceGt")
		priceGtFloat, err := strconv.ParseFloat(priceGt, 64)

		// Filtrar los productos que tengan un precio mayor a priceGt
		filteredProducts := make([]models.Product, 0)
		for _, product := range products {
			if product.Price > priceGtFloat {
				filteredProducts = append(filteredProducts, product)
			}
		}

		// Convertir los productos filtrados a JSON y enviar la respuesta
		response, err := json.Marshal(filteredProducts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

	http.ListenAndServe(":8080", r)
}
