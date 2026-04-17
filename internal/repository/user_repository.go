package repository

import (
	"agrocontrol-api/internal/domain/entities"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *entities.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
