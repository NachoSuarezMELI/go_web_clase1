package product

import "errors"

var (
	ErrProdNotFound     = errors.New("product not found")
	ErrProdInvalidField = errors.New("product is invalid")
)

type Product struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	CodeValue    string  `json:"code_value"`
	Is_Published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

type ResponseProducts struct {
	Message string    `json:"message"`
	Data    []Product `json:"data"`
}

type ResponseProduct struct {
	Message string  `json:"message"`
	Data    Product `json:"data"`
}

type RequestBodyProduct struct {
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	CodeValue    string  `json:"code_value"`
	Is_Published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

type ResponseBodyProduct struct {
	Message string  `json:"message"`
	Data    Product `json:"data"`
	Error   bool    `json:"error"`
}

type ProductRepository interface {
	GetAllProducts() ([]Product, error)
	GetProductById(id int) (*Product, error)
	CreateProduct(p *Product) error
	FindProductsByPriceGt(price float64) []Product
	UpdateOrCreateProduct(p *RequestBodyProduct, id int) error
	UpdatePartial(map[string]any, int) error
	DeleteProduct(id int) error
}

type ProductService interface {
	GetAllProducts() ([]Product, error)
	GetProductById(id int) (*Product, error)
	CreateProduct(p *Product) (err error)
	FindProductsByPriceGt(price float64) []Product
	UpdateOrCreateProduct(p *RequestBodyProduct, id int) error
	UpdatePartial(map[string]any, int) error
	DeleteProduct(id int) error
}
