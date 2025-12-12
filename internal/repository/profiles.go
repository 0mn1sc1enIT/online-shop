package repository

import (
	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ProfilesRepo struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewProfilesRepo(db *gorm.DB, logger *zerolog.Logger) *ProfilesRepo {
	return &ProfilesRepo{
		db:     db,
		logger: logger,
	}
}

func (r *ProfilesRepo) Create(profile *domain.Profile) error {
	if err := r.db.Create(profile).Error; err != nil {
		r.logger.Error().
			Err(err).
			Uint("user_id", profile.UserID).
			Str("first_name", profile.FirstName).
			Msg("Failed to create user profile")
		return err
	}

	r.logger.Debug().
		Uint("profile_id", profile.ID).
		Uint("user_id", profile.UserID).
		Msg("Profile created successfully")
	return nil
}

func (r *ProfilesRepo) Update(profile *domain.Profile) error {
	// Сохраняем все поля
	if err := r.db.Save(profile).Error; err != nil {
		r.logger.Error().
			Err(err).
			Uint("profile_id", profile.ID).
			Uint("user_id", profile.UserID).
			Msg("Failed to update profile")
		return err
	}

	r.logger.Debug().
		Uint("profile_id", profile.ID).
		Msg("Profile updated successfully")
	return nil
}

func (r *ProfilesRepo) GetByUserID(userID uint) (*domain.Profile, error) {
	var profile domain.Profile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		r.logger.Error().
			Err(err).
			Uint("user_id", userID).
			Msg("Failed to find profile by UserID")
		return nil, err
	}

	return &profile, nil
}
