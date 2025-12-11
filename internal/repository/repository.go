package repository

import (
	"github.com/OmniscienIT/GolangAPI/internal/domain"

	"gorm.io/gorm"
)

// Интерфейсы для каждой сущности
type Users interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	GetByID(id uint) (*domain.User, error)
}

type Profiles interface {
	Create(profile *domain.Profile) error
	Update(profile *domain.Profile) error
	GetByUserID(userID uint) (*domain.Profile, error)
}

type Categories interface {
	Create(category *domain.Category) error
	GetAll() ([]domain.Category, error)
	GetByID(id uint) (*domain.Category, error)
	Update(category *domain.Category) error
	Delete(id uint) error
}

type Products interface {
	Create(product *domain.Product) error
	GetAll() ([]domain.Product, error)
	GetByID(id uint) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uint) error
}

type Orders interface {
	Create(order *domain.Order) error
	GetAllByUserID(userID uint) ([]domain.Order, error)
	GetByID(id uint) (*domain.Order, error)
	UpdateStatus(id uint, status string) error
}

// Repositories - структура, собирающая все репозитории
type Repositories struct {
	Users      Users
	Profiles   Profiles
	Categories Categories
	Products   Products
	Orders     Orders
}

// Конструктор
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Users:      NewUsersRepo(db),
		Profiles:   NewProfilesRepo(db),
		Categories: NewCategoriesRepo(db),
		Products:   NewProductsRepo(db),
		Orders:     NewOrdersRepo(db),
	}
}
