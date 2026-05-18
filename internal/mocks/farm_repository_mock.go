package mocks

import (
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/domain/ports"
)

// FarmRepositoryMock implementa ports.FarmRepository para uso em testes.
// Cada campo Fn pode ser sobrescrito pelo teste para controlar o comportamento.
type FarmRepositoryMock struct {
	CreateFn                   func(farm *entities.Farm) error
	FindByIDFn                 func(id uint) (*entities.Farm, error)
	FindAllPaginatedFn         func(offset, limit int, search string) ([]entities.Farm, int64, error)
	FindByCreatedByPaginatedFn func(userID uint, offset, limit int, search string) ([]entities.Farm, int64, error)
	FindAllFn                  func() ([]entities.Farm, error)
	FindByCreatedByFn          func(userID uint) ([]entities.Farm, error)
	UpdateFn                   func(farm *entities.Farm) error
	DeleteFn                   func(id uint) error
	CountFn                    func() (int64, error)
}

// Garante que FarmRepositoryMock implementa a interface
var _ ports.FarmRepository = (*FarmRepositoryMock)(nil)

func (m *FarmRepositoryMock) Create(farm *entities.Farm) error {
	return m.CreateFn(farm)
}
func (m *FarmRepositoryMock) FindByID(id uint) (*entities.Farm, error) {
	return m.FindByIDFn(id)
}
func (m *FarmRepositoryMock) FindAllPaginated(offset, limit int, search string) ([]entities.Farm, int64, error) {
	return m.FindAllPaginatedFn(offset, limit, search)
}
func (m *FarmRepositoryMock) FindByCreatedByPaginated(userID uint, offset, limit int, search string) ([]entities.Farm, int64, error) {
	return m.FindByCreatedByPaginatedFn(userID, offset, limit, search)
}
func (m *FarmRepositoryMock) FindAll() ([]entities.Farm, error) {
	return m.FindAllFn()
}
func (m *FarmRepositoryMock) FindByCreatedBy(userID uint) ([]entities.Farm, error) {
	return m.FindByCreatedByFn(userID)
}
func (m *FarmRepositoryMock) Update(farm *entities.Farm) error {
	return m.UpdateFn(farm)
}
func (m *FarmRepositoryMock) Delete(id uint) error {
	return m.DeleteFn(id)
}
func (m *FarmRepositoryMock) Count() (int64, error) {
	return m.CountFn()
}
