package service

import (
	"web/clase1/internal"
	"web/clase1/internal/repository"
)

type ProductService struct {
	repository *repository.ProductRepository
}

func NewProductService(repository *repository.ProductRepository) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (s *ProductService) GetAllProducts() []product.Product {
	return s.repository.GetAllProducts()
}

func (s *ProductService) GetProductById(id int) (*product.Product, error) {
	p, err := s.repository.GetProductById(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProductService) CreateProduct(product *product.Product) (err error) {
	return s.repository.CreateProduct(product)
}

func (s *ProductService) FindProductsByPriceGt(price float64) []product.Product {
	return s.repository.FindProductsByPriceGt(price)
}
