package product

import (
	"encoding/json"
	"errors"
	"time"
)

func JsonHandler(data []byte) ([]Product, error) {
	var products []Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func GetProducts(st map[int]*Product) []Product {
	products := make([]Product, 0, len(st))
	for _, product := range st {
		products = append(products, *product)
	}
	return products
}

func FindProductById(id int, products map[int]*Product) *Product {
	product, ok := products[id]
	if !ok {
		return nil
	}
	return product
}
func FindProductsByPriceGt(price float64, products map[int]*Product) []Product {
	var productsFound []Product
	for _, product := range products {
		if product.Price > price {
			productsFound = append(productsFound, *product)
		}
	}
	return productsFound
}

func CreateProduct(id int, body RequestBodyProduct, products map[int]*Product) (*Product, error) {
	_, err := ValidateProductFields(body)
	if err != nil {
		return nil, err
	}
	_, err = ValidateProductCodeValue(body.CodeValue, products)
	if err != nil {
		return nil, err
	}

	_, err = ValidateExpirationDate(body.Expiration)
	if err != nil {
		return nil, err
	}

	product := Product{
		Id:           id,
		Name:         body.Name,
		Quantity:     body.Quantity,
		CodeValue:    body.CodeValue,
		Is_Published: body.Is_Published,
		Expiration:   body.Expiration,
		Price:        body.Price,
	}
	products[id] = &product
	return &product, nil
}

func ValidateProductFields(p RequestBodyProduct) (bool, error) {
	if p.Name == "" {
		return false, errors.New("invalid product, name is required")
	}
	if p.Quantity == 0 {
		return false, errors.New("invalid product, quantity is required")
	}
	if p.CodeValue == "" {
		return false, errors.New("invalid product, code value is required")
	}
	if p.Price == 0 {
		return false, errors.New("invalid product, price is required")
	}
	if p.Expiration == "" {
		return false, errors.New("invalid product, expiration is required")
	}
	return true, nil
}

func ValidateProductCodeValue(codeValue string, products map[int]*Product) (bool, error) {
	for _, product := range products {
		if product.CodeValue == codeValue {
			return false, errors.New("invalid product, code value already exists")
		}
	}
	return true, nil
}

func ValidateExpirationDate(expiration string) (bool, error) {
	if len(expiration) != 10 {
		return false, errors.New("invalid expiration format")
	}
	if expiration[2] != '/' || expiration[5] != '/' {
		return false, errors.New("invalid expiration format")
	}
	date, err := time.Parse("02/01/2006", expiration)
	if err != nil {
		return false, errors.New(err.Error())
	}
	if date.Before(time.Now()) {
		return false, errors.New("The expiration date must be greater than the current date")
	}
	return true, nil
}
