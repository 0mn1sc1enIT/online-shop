package service

import (
	"time"

	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/OmniscienIT/GolangAPI/internal/repository"
	"github.com/rs/zerolog"
)

type SignUpInput struct {
	Email    string
	Password string
}

type SignInInput struct {
	Email    string
	Password string
}

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
	Update(id uint, category domain.Category) error
	Delete(id uint) error
}

type Orders interface {
	Create(userID uint, inputOrder domain.Order) error
	GetAllByUserID(userID uint) ([]domain.Order, error)
}

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
	Logger    *zerolog.Logger
}

func NewServices(deps Deps) *Service {
	return &Service{
		Authorization: NewAuthService(deps.Repos.Users, deps.SignedKey, deps.TokenTTL, deps.Logger),
		Products:      NewProductsService(deps.Repos.Products, deps.Repos.Categories, deps.Logger),
		Categories:    NewCategoriesService(deps.Repos.Categories, deps.Logger),
		Orders:        NewOrdersService(deps.Repos.Orders, deps.Repos.Products, deps.Logger),
	}
}
