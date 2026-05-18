package mocks

import (
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/domain/ports"
)

type SeasonRepositoryMock struct {
	CreateFn           func(season *entities.Season) error
	FindByIDFn         func(id uint) (*entities.Season, error)
	FindAllPaginatedFn func(offset, limit int, search string) ([]entities.Season, int64, error)
	FindAllFn          func() ([]entities.Season, error)
	UpdateFn           func(season *entities.Season) error
	DeleteFn           func(id uint) error
}

var _ ports.SeasonRepository = (*SeasonRepositoryMock)(nil)

func (m *SeasonRepositoryMock) Create(s *entities.Season) error { return m.CreateFn(s) }
func (m *SeasonRepositoryMock) FindByID(id uint) (*entities.Season, error) { return m.FindByIDFn(id) }
func (m *SeasonRepositoryMock) FindAllPaginated(o, l int, search string) ([]entities.Season, int64, error) { return m.FindAllPaginatedFn(o, l, search) }
func (m *SeasonRepositoryMock) FindAll() ([]entities.Season, error) { return m.FindAllFn() }
func (m *SeasonRepositoryMock) Update(s *entities.Season) error { return m.UpdateFn(s) }
func (m *SeasonRepositoryMock) Delete(id uint) error { return m.DeleteFn(id) }
