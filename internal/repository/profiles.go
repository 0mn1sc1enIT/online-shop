package repository

import (
	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"gorm.io/gorm"
)

type ProfilesRepo struct {
	db *gorm.DB
}

func NewProfilesRepo(db *gorm.DB) *ProfilesRepo {
	return &ProfilesRepo{db: db}
}

func (r *ProfilesRepo) Create(profile *domain.Profile) error {
	return r.db.Create(profile).Error
}

func (r *ProfilesRepo) Update(profile *domain.Profile) error {
	// Сохраняем все поля
	return r.db.Save(profile).Error
}

func (r *ProfilesRepo) GetByUserID(userID uint) (*domain.Profile, error) {
	var profile domain.Profile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}
