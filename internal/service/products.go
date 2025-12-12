package service

import (
	"errors"

	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/OmniscienIT/GolangAPI/internal/repository"
	"github.com/rs/zerolog"
)

type ProductsService struct {
	repo    repository.Products
	catRepo repository.Categories
	logger  *zerolog.Logger
}

func NewProductsService(repo repository.Products, catRepo repository.Categories, logger *zerolog.Logger) *ProductsService {
	return &ProductsService{
		repo:    repo,
		catRepo: catRepo,
		logger:  logger,
	}
}

func (s *ProductsService) Create(product domain.Product) error {
	_, err := s.catRepo.GetByID(product.CategoryID)
	if err != nil {
		s.logger.Warn().
			Uint("category_id", product.CategoryID).
			Str("product_name", product.Name).
			Msg("Create product failed: category does not exist")
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
	_, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Warn().
			Uint("product_id", id).
			Err(err).
			Msg("Update product failed: product not found")
		return err
	}

	product.ID = id
	return s.repo.Update(&product)
}

func (s *ProductsService) Delete(id uint) error {
	return s.repo.Delete(id)
}
