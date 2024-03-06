package services

import (
	"encoding/json"
	"web/clase1/models"
)

func JsonHandler(data []byte) ([]models.Product, error) {
	var products []models.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func FindProductById(id int, products []models.Product) *models.Product {
	for _, product := range products {
		if product.Id == id {
			return &product
		}
	}
	return nil
}
