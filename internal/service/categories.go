package service

import (
	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/OmniscienIT/GolangAPI/internal/repository"
)

type CategoriesService struct {
	repo repository.Categories
}

func NewCategoriesService(repo repository.Categories) *CategoriesService {
	return &CategoriesService{repo: repo}
}

func (s *CategoriesService) Create(category domain.Category) error {
	return s.repo.Create(&category)
}

func (s *CategoriesService) GetAll() ([]domain.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoriesService) GetByID(id uint) (domain.Category, error) {
	cat, err := s.repo.GetByID(id)
	if err != nil {
		return domain.Category{}, err
	}
	return *cat, nil
}
