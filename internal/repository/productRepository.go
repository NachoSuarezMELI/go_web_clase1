package repository

import (
	"encoding/json"
	"errors"
	"web/clase1/internal"
	"web/clase1/internal/storage"
)

type ProductSlice struct {
	slice   []product.Product
	storage storage.Storage
}

func NewProductRepository(st storage.Storage) *ProductSlice {
	data, err := st.Read()
	if err != nil {
		return nil
	}

	//convert data to slice of products
	var products []product.Product
	err = json.Unmarshal(data, &products)
	if err != nil {
		return nil
	}

	return &ProductSlice{
		slice:   products,
		storage: st,
	}
}

func (r *ProductSlice) GetAllProducts() ([]product.Product, error) {
	if len(r.slice) == 0 {
		return nil, errors.New("no products found")
	}
	return r.slice, nil
}

func (r *ProductSlice) GetProductById(id int) (*product.Product, error) {
	for _, product := range r.slice {
		if product.Id == id {
			return &product, nil
		}
	}
	return nil, product.ErrProdNotFound

}

func (r *ProductSlice) FindProductsByPriceGt(price float64) []product.Product {
	var productsFound []product.Product
	for _, product := range r.slice {
		if product.Price > price {
			productsFound = append(productsFound, product)
		}
	}
	return productsFound
}

func (r *ProductSlice) CreateProduct(p *product.Product) (err error) {

	p.Id = len(r.slice) + 1
	r.slice = append(r.slice, *p)

	//save slice to storage
	data, err := json.Marshal(r.slice)
	if err != nil {
		return err
	}
	err = r.storage.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductSlice) UpdateOrCreateProduct(p *product.RequestBodyProduct, id int) error {
	_, err := r.GetProductById(id)
	if err != nil {
		newProduct := product.Product{
			Id:           len(r.slice) + 1,
			Name:         p.Name,
			Quantity:     p.Quantity,
			CodeValue:    p.CodeValue,
			Is_Published: p.Is_Published,
			Expiration:   p.Expiration,
			Price:        p.Price,
		}
		r.slice = append(r.slice, *&newProduct)
	} else {
		product := r.slice[(id - 1)]
		product.Name = p.Name
		product.Quantity = p.Quantity
		product.CodeValue = p.CodeValue
		product.Is_Published = p.Is_Published
		product.Expiration = p.Expiration
		product.Price = p.Price
		r.slice[(id - 1)] = product
	}

	//save slice to storage
	data, err := json.Marshal(r.slice)
	if err != nil {
		return err
	}
	err = r.storage.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePartial updates a product by id
func (r *ProductSlice) UpdatePartial(p map[string]interface{}, id int) error {
	product, err := r.GetProductById(id)
	if err != nil {
		return err
	}
	for key, value := range p {
		switch key {
		case "name":
			product.Name = value.(string)
		case "quantity":
			product.Quantity = value.(int)
		case "code_value":
			product.CodeValue = value.(string)
		case "is_published":
			product.Is_Published = value.(bool)
		case "expiration":
			product.Expiration = value.(string)
		case "price":
			product.Price = value.(float64)
		default:
			return errors.New("invalid field")
		}
	}
	r.slice[(id - 1)] = *product

	//save slice to storage
	data, err := json.Marshal(r.slice)
	if err != nil {
		return err
	}
	err = r.storage.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductSlice) DeleteProduct(id int) error {
	_, err := r.GetProductById(id)
	if err != nil {
		return product.ErrProdNotFound
	}

	// Delete product
	for i, product := range r.slice {
		if product.Id == id {
			r.slice = append(r.slice[:i], r.slice[i+1:]...)
			break
		}
	}

	// Save slice to storage
	data, err := json.Marshal(r.slice)
	if err != nil {
		return err
	}
	err = r.storage.Write(data)
	if err != nil {
		return err
	}
	return nil

}
