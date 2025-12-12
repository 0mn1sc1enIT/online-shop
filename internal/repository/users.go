package repository

import (
	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"gorm.io/gorm"
)

type UsersRepo struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UsersRepo) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Profile").Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UsersRepo) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Profile").First(&user, id).Error
	return &user, err
}
