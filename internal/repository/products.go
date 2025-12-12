package repository

import (
	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"gorm.io/gorm"
)

type ProductsRepo struct {
	db *gorm.DB
}

func NewProductsRepo(db *gorm.DB) *ProductsRepo {
	return &ProductsRepo{db: db}
}

func (r *ProductsRepo) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductsRepo) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Preload("Category").Find(&products).Error
	return products, err
}

func (r *ProductsRepo) GetByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Category").First(&product, id).Error
	return &product, err
}

func (r *ProductsRepo) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductsRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}
