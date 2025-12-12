package service

import (
	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/OmniscienIT/GolangAPI/internal/repository"
	"github.com/rs/zerolog"
)

type CategoriesService struct {
	repo   repository.Categories
	logger *zerolog.Logger
}

func NewCategoriesService(repo repository.Categories, logger *zerolog.Logger) *CategoriesService {
	return &CategoriesService{
		repo:   repo,
		logger: logger,
	}
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

func (s *CategoriesService) Update(id uint, category domain.Category) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Warn().
			Uint("category_id", id).
			Err(err).
			Msg("Update category failed: not found")
		return err
	}

	category.ID = id
	return s.repo.Update(&category)
}

func (s *CategoriesService) Delete(id uint) error {
	return s.repo.Delete(id)
}
