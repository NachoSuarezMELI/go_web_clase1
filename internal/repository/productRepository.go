package repository

import (
	"encoding/json"
	"errors"
	"web/clase1/internal"
	"web/clase1/platform/tools"
)

type ProductRepository struct {
	storage []product.Product
	lastId  int
}

func NewProductRepository(db []product.Product, lastId int) *ProductRepository {

	return &ProductRepository{
		storage: db,
		lastId:  lastId,
	}
}

func (r *ProductRepository) GetAllProducts() []product.Product {
	products := make([]product.Product, 0)
	for _, product := range r.storage {
		products = append(products, product)
	}
	return products
}

func (r *ProductRepository) GetProductById(id int) (*product.Product, error) {
	for _, product := range r.storage {
		if product.Id == id {
			return &product, nil
		}
	}
	return nil, errors.New("product not found")

}

func (r *ProductRepository) CreateProduct(p *product.Product) (err error) {
	r.lastId++
	p.Id = r.lastId
	r.storage[p.Id] = *p
	return nil
}

func (r *ProductRepository) FindProductsByPriceGt(price float64) []product.Product {
	var productsFound []product.Product
	for _, product := range r.storage {
		if product.Price > price {
			productsFound = append(productsFound, product)
		}
	}
	return productsFound
}

func (r *ProductRepository) GeneretaDB() error {
	bytes, err := tools.ReadFile("../docs/db/products.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &r.storage); err != nil {
		return err
	}
	return nil
}
