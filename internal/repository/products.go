package repository

import (
	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ProductsRepo struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewProductsRepo(db *gorm.DB, logger *zerolog.Logger) *ProductsRepo {
	return &ProductsRepo{
		db:     db,
		logger: logger,
	}
}

func (r *ProductsRepo) Create(product *domain.Product) error {
	if err := r.db.Create(product).Error; err != nil {
		r.logger.Error().
			Err(err).
			Str("product_name", product.Name).
			Float64("price", product.Price).
			Uint("category_id", product.CategoryID).
			Msg("Failed to create product in DB")
		return err
	}

	r.logger.Debug().
		Uint("product_id", product.ID).
		Msg("Product created successfully")
	return nil
}

func (r *ProductsRepo) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	if err := r.db.Preload("Category").Find(&products).Error; err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to fetch all products")
		return nil, err
	}
	return products, nil
}

func (r *ProductsRepo) GetByID(id uint) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.Preload("Category").First(&product, id).Error; err != nil {
		r.logger.Error().
			Err(err).
			Uint("product_id", id).
			Msg("Failed to get product by ID")
		return nil, err
	}
	return &product, nil
}

func (r *ProductsRepo) Update(product *domain.Product) error {
	if err := r.db.Save(product).Error; err != nil {
		r.logger.Error().
			Err(err).
			Uint("product_id", product.ID).
			Msg("Failed to update product")
		return err
	}

	r.logger.Debug().
		Uint("product_id", product.ID).
		Msg("Product updated successfully")
	return nil
}

func (r *ProductsRepo) Delete(id uint) error {
	if err := r.db.Delete(&domain.Product{}, id).Error; err != nil {
		r.logger.Error().
			Err(err).
			Uint("product_id", id).
			Msg("Failed to delete product")
		return err
	}

	r.logger.Debug().
		Uint("product_id", id).
		Msg("Product deleted")
	return nil
}
