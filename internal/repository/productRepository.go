package repository

import (
	"encoding/json"
	"errors"
	"web/clase1/internal"
	"web/clase1/platform/tools"
)

type ProductSlice struct {
	storage []product.Product
	lastId  int
}

func NewProductRepository(data []product.Product, lastId int) *ProductSlice {
	if data == nil {
		data = make([]product.Product, 0)
	}

	return &ProductSlice{
		storage: data,
		lastId:  lastId,
	}
}

func (r *ProductSlice) GetAllProducts() ([]product.Product, error) {
	if len(r.storage) == 0 {
		return nil, errors.New("no products found")
	}
	return r.storage, nil
}

func (r *ProductSlice) GetProductById(id int) (*product.Product, error) {
	for _, product := range r.storage {
		if product.Id == id {
			return &product, nil
		}
	}
	return nil, errors.New("product not found")

}

func (r *ProductSlice) FindProductsByPriceGt(price float64) []product.Product {
	var productsFound []product.Product
	for _, product := range r.storage {
		if product.Price > price {
			productsFound = append(productsFound, product)
		}
	}
	return productsFound
}

func (r *ProductSlice) CreateProduct(p *product.Product) (err error) {
	r.lastId++
	p.Id = r.lastId
	r.storage = append(r.storage, *p)
	return nil
}

func (r *ProductSlice) UpdateOrCreateProduct(p *product.RequestBodyProduct, id int) error {
	_, err := r.GetProductById(id)
	if err != nil {
		r.lastId++
		newProduct := product.Product{
			Id:           r.lastId,
			Name:         p.Name,
			Quantity:     p.Quantity,
			CodeValue:    p.CodeValue,
			Is_Published: p.Is_Published,
			Expiration:   p.Expiration,
			Price:        p.Price,
		}
		r.storage = append(r.storage, *&newProduct)
	} else {
		product := r.storage[(id - 1)]
		product.Name = p.Name
		product.Quantity = p.Quantity
		product.CodeValue = p.CodeValue
		product.Is_Published = p.Is_Published
		product.Expiration = p.Expiration
		product.Price = p.Price
		r.storage[(id - 1)] = product
	}
	return nil
}

// // UpdateSomeAtributes updates some atributes of a product

// // TODO: implement this method

// GeneretaDB generates the database from a json file
func (r *ProductSlice) GeneretaDB() error {
	bytes, err := tools.ReadFile("../docs/db/products.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &r.storage); err != nil {
		return err
	}
	r.lastId = len(r.storage)
	return nil
}
