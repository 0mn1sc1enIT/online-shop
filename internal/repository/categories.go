package repository

import (
	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"gorm.io/gorm"
)

type CategoriesRepo struct {
	db *gorm.DB
}

func NewCategoriesRepo(db *gorm.DB) *CategoriesRepo {
	return &CategoriesRepo{db: db}
}

func (r *CategoriesRepo) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoriesRepo) GetAll() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *CategoriesRepo) GetByID(id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	return &category, err
}

func (r *CategoriesRepo) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *CategoriesRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}
