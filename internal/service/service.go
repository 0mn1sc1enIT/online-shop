package service

import (
	"time"

	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/OmniscienIT/GolangAPI/internal/repository"
)

// Структура для входных данных авторизации
type SignUpInput struct {
	Email    string
	Password string
}

type SignInInput struct {
	Email    string
	Password string
}

// Интерфейсы
type Authorization interface {
	CreateUser(user domain.User) (uint, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (uint, string, error) // возвращает userID и role
}

type Products interface {
	Create(product domain.Product) error
	GetAll() ([]domain.Product, error)
	GetByID(id uint) (domain.Product, error)
	Update(id uint, product domain.Product) error
	Delete(id uint) error
}

type Categories interface {
	Create(category domain.Category) error
	GetAll() ([]domain.Category, error)
	GetByID(id uint) (domain.Category, error)
}

type Orders interface {
	Create(userID uint, inputOrder domain.Order) error // inputOrder содержит список Items
	GetAllByUserID(userID uint) ([]domain.Order, error)
}

// Service собирает все интерфейсы
type Service struct {
	Authorization
	Products
	Categories
	Orders
}

type Deps struct {
	Repos     *repository.Repositories
	TokenTTL  time.Duration
	SignedKey string
}

// Конструктор сервисов
func NewServices(deps Deps) *Service {
	return &Service{
		Authorization: NewAuthService(deps.Repos.Users, deps.SignedKey, deps.TokenTTL),
		Products:      NewProductsService(deps.Repos.Products, deps.Repos.Categories),
		Categories:    NewCategoriesService(deps.Repos.Categories),
		Orders:        NewOrdersService(deps.Repos.Orders, deps.Repos.Products),
	}
}
