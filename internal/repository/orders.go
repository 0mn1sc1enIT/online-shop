package repository

import (
	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"gorm.io/gorm"
)

type OrdersRepo struct {
	db *gorm.DB
}

func NewOrdersRepo(db *gorm.DB) *OrdersRepo {
	return &OrdersRepo{db: db}
}

func (r *OrdersRepo) Create(order *domain.Order) error {
	// GORM автоматически создаст записи в OrderItems, так как они вложены в Order
	return r.db.Create(order).Error
}

func (r *OrdersRepo) GetAllByUserID(userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (r *OrdersRepo) GetByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items.Product").First(&order, id).Error
	return &order, err
}

func (r *OrdersRepo) UpdateStatus(id uint, status string) error {
	return r.db.Model(&domain.Order{}).Where("id = ?", id).Update("status", status).Error
}
