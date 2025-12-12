package repository

import (
	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type OrdersRepo struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewOrdersRepo(db *gorm.DB, logger *zerolog.Logger) *OrdersRepo {
	return &OrdersRepo{
		db:     db,
		logger: logger,
	}
}

func (r *OrdersRepo) Create(order *domain.Order) error {
	if err := r.db.Create(order).Error; err != nil {
		r.logger.Error().
			Err(err).
			Uint("user_id", order.UserID).
			Float64("total_price", order.TotalPrice).
			Msg("Failed to create order in DB")
		return err
	}

	r.logger.Debug().
		Uint("order_id", order.ID).
		Int("items_count", len(order.Items)).
		Msg("Order created successfully")
	return nil
}

func (r *OrdersRepo) GetAllByUserID(userID uint) ([]domain.Order, error) {
	var orders []domain.Order

	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		r.logger.Error().
			Err(err).
			Uint("user_id", userID).
			Msg("Failed to fetch orders for user")
		return nil, err
	}

	return orders, nil
}

func (r *OrdersRepo) GetByID(id uint) (*domain.Order, error) {
	var order domain.Order

	err := r.db.Preload("Items.Product").First(&order, id).Error
	if err != nil {
		r.logger.Error().
			Err(err).
			Uint("order_id", id).
			Msg("Failed to fetch order by ID")
		return nil, err
	}

	return &order, nil
}

func (r *OrdersRepo) UpdateStatus(id uint, status string) error {
	err := r.db.Model(&domain.Order{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		r.logger.Error().
			Err(err).
			Uint("order_id", id).
			Str("new_status", status).
			Msg("Failed to update order status")
		return err
	}

	r.logger.Debug().
		Uint("order_id", id).
		Str("status", status).
		Msg("Order status updated")
	return nil
}
