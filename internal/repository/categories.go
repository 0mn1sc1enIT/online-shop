package repository

import (
	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CategoriesRepo struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewCategoriesRepo(db *gorm.DB, logger *zerolog.Logger) *CategoriesRepo {
	return &CategoriesRepo{
		db:     db,
		logger: logger,
	}
}

func (r *CategoriesRepo) Create(category *domain.Category) error {
	if err := r.db.Create(category).Error; err != nil {
		r.logger.Error().
			Err(err).
			Str("category_name", category.Name).
			Msg("Failed to create category in DB")
		return err
	}

	r.logger.Debug().
		Uint("category_id", category.ID).
		Msg("Category created successfully")
	return nil
}

func (r *CategoriesRepo) GetAll() ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.db.Find(&categories).Error; err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to fetch categories")
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepo) GetByID(id uint) (*domain.Category, error) {
	var category domain.Category
	if err := r.db.First(&category, id).Error; err != nil {
		r.logger.Error().
			Err(err).
			Uint("category_id", id).
			Msg("Failed to get category by ID")
		return nil, err
	}
	return &category, nil
}

func (r *CategoriesRepo) Update(category *domain.Category) error {
	if err := r.db.Save(category).Error; err != nil {
		r.logger.Error().
			Err(err).
			Uint("category_id", category.ID).
			Msg("Failed to update category")
		return err
	}

	r.logger.Debug().
		Uint("category_id", category.ID).
		Msg("Category updated successfully")
	return nil
}

func (r *CategoriesRepo) Delete(id uint) error {
	if err := r.db.Delete(&domain.Category{}, id).Error; err != nil {
		r.logger.Error().
			Err(err).
			Uint("category_id", id).
			Msg("Failed to delete category")
		return err
	}

	r.logger.Debug().
		Uint("category_id", id).
		Msg("Category deleted")
	return nil
}
