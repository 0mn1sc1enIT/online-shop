package domain

import (
	"gorm.io/gorm"
)

// Пользователь системы
type User struct {
	gorm.Model
	Email    string  `gorm:"uniqueIndex;not null" json:"email"`
	Password string  `gorm:"not null" json:"-"`
	Role     string  `gorm:"default:'user'" json:"role"`
	Profile  Profile `json:"profile,omitempty"`
	Orders   []Order `json:"orders,omitempty"`
}

// Информация о пользователе (1 к 1 с User)
type Profile struct {
	gorm.Model
	UserID    uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
}

// Категория товара (1 ко многим с Product)
type Category struct {
	gorm.Model
	Name     string    `gorm:"unique;not null" json:"name"`
	Products []Product `json:"products,omitempty"`
}

// Товар
type Product struct {
	gorm.Model
	Name        string   `gorm:"not null" json:"name"`
	Description string   `json:"description"`
	Price       float64  `gorm:"not null" json:"price"`
	Stock       int      `gorm:"not null;default:0" json:"stock"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `json:"category,omitempty"`
}

// Заказ пользователя
type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id"`
	User       User        `json:"-"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `gorm:"default:'pending'" json:"status"` // pending, paid, shipped, cancelled
	Items      []OrderItem `json:"items"`
}

// Позиция в заказе (промежуточная таблица для M:M)
type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // Цена на момент покупки (цена товара может изменится)
}
