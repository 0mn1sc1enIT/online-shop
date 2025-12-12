package repository

import (
	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UsersRepo struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewUsersRepo(db *gorm.DB, logger *zerolog.Logger) *UsersRepo {
	return &UsersRepo{
		db:     db,
		logger: logger,
	}
}

func (r *UsersRepo) Create(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		r.logger.Error().
			Err(err).
			Str("email", user.Email).
			Msg("Failed to create user")
		return err
	}

	r.logger.Debug().
		Uint("user_id", user.ID).
		Str("email", user.Email).
		Msg("User created successfully")
	return nil
}

func (r *UsersRepo) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Profile").Where("email = ?", email).First(&user).Error
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("email", email).
			Msg("Failed to find user by email")
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepo) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Profile").First(&user, id).Error
	if err != nil {
		r.logger.Error().
			Err(err).
			Uint("user_id", id).
			Msg("Failed to find user by ID")
		return nil, err
	}

	return &user, nil
}
