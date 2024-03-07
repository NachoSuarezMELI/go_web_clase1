package services

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"
	"web/clase1/models"
)

func JsonHandler(data []byte) ([]models.Product, error) {
	var products []models.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func GetProducts(st map[int]*models.Product) []models.Product {
	products := make([]models.Product, 0, len(st))
	for _, product := range st {
		products = append(products, *product)
	}
	return products
}

func FindProductsByPriceGt(price float64, products map[int]*models.Product) []models.Product {
	var productsFound []models.Product
	for _, product := range products {
		if product.Price > price {
			productsFound = append(productsFound, *product)
		}
	}
	return productsFound
}

func FindProductById(id int, products map[int]*models.Product) *models.Product {
	product, ok := products[id]
	if !ok {
		return nil
	}
	return product
}

func CreateProduct(id int, body models.RequestBodyProduct, products map[int]*models.Product) (*models.Product, error) {
	_, err := ValidateProductFields(body)
	if err != nil {
		return nil, err
	}
	_, err = ValidateProductCodeValue(body.CodeValue, products)
	if err != nil {
		return nil, err
	}
	_, err = ValidateAllTheTypes(body)
	if err != nil {
		return nil, err
	}

	_, err = ValidateExpirationDate(body.Expiration)
	if err != nil {
		return nil, err
	}

	product := models.Product{
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

func ValidateProductFields(product models.RequestBodyProduct) (bool, error) {
	if product.Name == "" || product.Quantity == 0 || product.CodeValue == "" || product.Expiration == "" || product.Price == 0 {
		return false, errors.New("invalid product, some fields are empty")
	}
	if product.Price < 0 {
		return false, errors.New("invalid product, price must be greater than zero")
	}
	if product.Quantity < 0 {
		return false, errors.New("invalid product, quantity must be greater than zero")
	}
	return true, nil
}

func ValidateProductCodeValue(codeValue string, products map[int]*models.Product) (bool, error) {
	for _, product := range products {
		if product.CodeValue == codeValue {
			return false, errors.New("invalid product, code value already exists")
		}
	}
	return true, nil
}

func ValidateType(value interface{}, typeValue string) (bool, error) {
	if reflect.TypeOf(value).String() != typeValue {
		return false, errors.New("invalid type")
	}
	return true, nil
}

func ValidateAllTheTypes(values models.RequestBodyProduct) (bool, error) {
	if _, err := ValidateType(values.Name, "string"); err != nil {
		return false, errors.New("invalid type")
	}
	if _, err := ValidateType(values.Quantity, "int"); err != nil {
		return false, errors.New("invalid type")
	}
	if _, err := ValidateType(values.CodeValue, "string"); err != nil {
		return false, errors.New("invalid type")
	}
	if _, err := ValidateType(values.Is_Published, "bool"); err != nil {
		return false, errors.New("invalid type")
	}
	if _, err := ValidateType(values.Expiration, "string"); err != nil {
		return false, errors.New("invalid type")
	}
	if _, err := ValidateType(values.Price, "float64"); err != nil {
		return false, errors.New("invalid type")
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
