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

func (r *UserRepository) FindByID(id uint) (*entities.User, error) {
	var user entities.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAll() ([]entities.User, error) {
	var users []entities.User
	err := r.DB.Order("created_at desc").Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(user *entities.User) error {
	return r.DB.Save(user).Error
}
