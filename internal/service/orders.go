package service

import (
	"errors"

	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/OmniscienIT/GolangAPI/internal/repository"
)

type OrdersService struct {
	repo        repository.Orders
	productRepo repository.Products
}

func NewOrdersService(repo repository.Orders, productRepo repository.Products) *OrdersService {
	return &OrdersService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *OrdersService) Create(userID uint, order domain.Order) error {
	var totalPrice float64
	var orderItems []domain.OrderItem

	// Проходимся по товарам, которые хочет купить пользователь
	for _, item := range order.Items {
		// Получаем актуальную цену товара из БД
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return errors.New("product not found")
		}

		// Проверяем наличие на складе (опционально)
		if product.Stock < item.Quantity {
			return errors.New("not enough stock for product: " + product.Name)
		}

		// Формируем позицию заказа с актуальной ценой
		orderItem := domain.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		// Считаем сумму
		totalPrice += product.Price * float64(item.Quantity)
		orderItems = append(orderItems, orderItem)

		// Уменьшаем Stock у товара
		product.Stock -= item.Quantity
		err = s.productRepo.Update(product)
		if err != nil {
			return err
		}
	}

	if totalPrice == 0 {
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
