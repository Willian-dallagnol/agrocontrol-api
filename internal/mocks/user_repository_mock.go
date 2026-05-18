package mocks

import (
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/domain/ports"
)

type UserRepositoryMock struct {
	CreateFn      func(user *entities.User) error
	FindByIDFn    func(id uint) (*entities.User, error)
	FindByEmailFn func(email string) (*entities.User, error)
	FindAllFn     func() ([]entities.User, error)
	UpdateFn      func(user *entities.User) error
}

var _ ports.UserRepository = (*UserRepositoryMock)(nil)

func (m *UserRepositoryMock) Create(user *entities.User) error          { return m.CreateFn(user) }
func (m *UserRepositoryMock) FindByID(id uint) (*entities.User, error)  { return m.FindByIDFn(id) }
func (m *UserRepositoryMock) FindByEmail(email string) (*entities.User, error) { return m.FindByEmailFn(email) }
func (m *UserRepositoryMock) FindAll() ([]entities.User, error)         { return m.FindAllFn() }
func (m *UserRepositoryMock) Update(user *entities.User) error          { return m.UpdateFn(user) }
