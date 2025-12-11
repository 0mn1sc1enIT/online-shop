package service

import (
	"errors"

	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/OmniscienIT/GolangAPI/internal/repository"
)

type ProductsService struct {
	repo    repository.Products
	catRepo repository.Categories
}

func NewProductsService(repo repository.Products, catRepo repository.Categories) *ProductsService {
	return &ProductsService{repo: repo, catRepo: catRepo}
}

func (s *ProductsService) Create(product domain.Product) error {
	// Проверяем, существует ли категория
	_, err := s.catRepo.GetByID(product.CategoryID)
	if err != nil {
		return errors.New("category does not exist")
	}
	return s.repo.Create(&product)
}

func (s *ProductsService) GetAll() ([]domain.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductsService) GetByID(id uint) (domain.Product, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return domain.Product{}, err
	}
	return *p, nil
}

func (s *ProductsService) Update(id uint, product domain.Product) error {
	// Проверяем существование
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	product.ID = id
	return s.repo.Update(&product)
}

func (s *ProductsService) Delete(id uint) error {
	return s.repo.Delete(id)
}
