package mocks

import (
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/domain/ports"
)

type FieldRepositoryMock struct {
	CreateFn               func(field *entities.Field) error
	FindByIDFn             func(id uint) (*entities.Field, error)
	FindAllPaginatedFn     func(offset, limit int, search string) ([]entities.Field, int64, error)
	FindByUserPaginatedFn  func(userID uint, offset, limit int, search string) ([]entities.Field, int64, error)
	FindAllFn              func() ([]entities.Field, error)
	FindByUserFn           func(userID uint) ([]entities.Field, error)
	FindByFarmIDFn         func(farmID uint) ([]entities.Field, error)
	ExistsByNameAndFarmFn  func(name string, farmID uint, excludeID uint) (bool, error)
	UpdateFn               func(field *entities.Field) error
	DeleteFn               func(id uint) error
	CountFn                func() (int64, error)
	BelongsToUserFn        func(fieldID, userID uint) (bool, error)
}

var _ ports.FieldRepository = (*FieldRepositoryMock)(nil)

func (m *FieldRepositoryMock) Create(field *entities.Field) error                  { return m.CreateFn(field) }
func (m *FieldRepositoryMock) FindByID(id uint) (*entities.Field, error)           { return m.FindByIDFn(id) }
func (m *FieldRepositoryMock) FindAllPaginated(o, l int, s string) ([]entities.Field, int64, error) { return m.FindAllPaginatedFn(o, l, s) }
func (m *FieldRepositoryMock) FindByUserPaginated(uID uint, o, l int, s string) ([]entities.Field, int64, error) { return m.FindByUserPaginatedFn(uID, o, l, s) }
func (m *FieldRepositoryMock) FindAll() ([]entities.Field, error)                  { return m.FindAllFn() }
func (m *FieldRepositoryMock) FindByUser(userID uint) ([]entities.Field, error)    { return m.FindByUserFn(userID) }
func (m *FieldRepositoryMock) FindByFarmID(farmID uint) ([]entities.Field, error)  { return m.FindByFarmIDFn(farmID) }
func (m *FieldRepositoryMock) ExistsByNameAndFarm(name string, farmID, exID uint) (bool, error) { return m.ExistsByNameAndFarmFn(name, farmID, exID) }
func (m *FieldRepositoryMock) Update(field *entities.Field) error                  { return m.UpdateFn(field) }
func (m *FieldRepositoryMock) Delete(id uint) error                                { return m.DeleteFn(id) }
func (m *FieldRepositoryMock) Count() (int64, error)                               { return m.CountFn() }
func (m *FieldRepositoryMock) BelongsToUser(fieldID, userID uint) (bool, error)   { return m.BelongsToUserFn(fieldID, userID) }
