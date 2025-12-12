package service

import (
	"errors"

	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/OmniscienIT/GolangAPI/internal/repository"
	"github.com/rs/zerolog"
)

type OrdersService struct {
	repo        repository.Orders
	productRepo repository.Products
	logger      *zerolog.Logger
}

func NewOrdersService(repo repository.Orders, productRepo repository.Products, logger *zerolog.Logger) *OrdersService {
	return &OrdersService{
		repo:        repo,
		productRepo: productRepo,
		logger:      logger,
	}
}

func (s *OrdersService) Create(userID uint, order domain.Order) error {
	var totalPrice float64
	var orderItems []domain.OrderItem

	for _, item := range order.Items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			s.logger.Warn().
				Uint("product_id", item.ProductID).
				Uint("user_id", userID).
				Msg("Order failed: product not found")
			return errors.New("product not found")
		}

		if product.Stock < item.Quantity {
			s.logger.Warn().
				Str("product_name", product.Name).
				Int("stock_available", product.Stock).
				Int("requested", item.Quantity).
				Msg("Order failed: not enough stock")
			return errors.New("not enough stock for product: " + product.Name)
		}

		orderItem := domain.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		totalPrice += product.Price * float64(item.Quantity)
		orderItems = append(orderItems, orderItem)

		product.Stock -= item.Quantity
		err = s.productRepo.Update(product)
		if err != nil {
			s.logger.Error().
				Err(err).
				Uint("product_id", product.ID).
				Msg("Failed to update product stock during order")
			return err
		}
	}

	if totalPrice == 0 {
		s.logger.Warn().
			Uint("user_id", userID).
			Msg("Order failed: order is empty or total price 0")
		return errors.New("order is empty")
	}

	newOrder := domain.Order{
		UserID:     userID,
		Items:      orderItems,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	return s.repo.Create(&newOrder)
}

func (s *OrdersService) GetAllByUserID(userID uint) ([]domain.Order, error) {
	return s.repo.GetAllByUserID(userID)
}
