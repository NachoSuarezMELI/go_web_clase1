package service

import (
	"web/clase1/internal"
)

type Service struct {
	repository product.ProductRepository
}

func NewProductService(repository product.ProductRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAllProducts() ([]product.Product, error) {
	return s.repository.GetAllProducts()
}

func (s *Service) GetProductById(id int) (*product.Product, error) {
	p, err := s.repository.GetProductById(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) FindProductsByPriceGt(price float64) []product.Product {
	return s.repository.FindProductsByPriceGt(price)
}

func (s *Service) CreateProduct(product *product.Product) (err error) {
	return s.repository.CreateProduct(product)
}

func (s *Service) UpdateOrCreateProduct(product *product.RequestBodyProduct, id int) error {
	return s.repository.UpdateOrCreateProduct(product, id)
}
