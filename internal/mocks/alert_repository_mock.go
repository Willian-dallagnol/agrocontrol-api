package mocks

import (
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/domain/ports"
)

type AlertRepositoryMock struct {
	FindByIDFn       func(id uint) (*entities.Alert, error)
	FindAllFn        func() ([]entities.Alert, error)
	FindByUserFn     func(userID uint) ([]entities.Alert, error)
	FindOpenFn       func() ([]entities.Alert, error)
	FindOpenByUserFn func(userID uint) ([]entities.Alert, error)
	UpdateFn         func(alert *entities.Alert) error
	DeleteFn         func(id uint) error
	CountOpenFn      func() (int64, error)
	CountOpenByUserFn func(userID uint) (int64, error)
	CreateTxFn       func(tx ports.TxRunner, alert *entities.Alert) error
	CreateFn         func(alert *entities.Alert) error
}

var _ ports.AlertRepository = (*AlertRepositoryMock)(nil)

func (m *AlertRepositoryMock) FindByID(id uint) (*entities.Alert, error)         { return m.FindByIDFn(id) }
func (m *AlertRepositoryMock) FindAll() ([]entities.Alert, error)                { return m.FindAllFn() }
func (m *AlertRepositoryMock) FindByUser(userID uint) ([]entities.Alert, error)  { return m.FindByUserFn(userID) }
func (m *AlertRepositoryMock) FindOpen() ([]entities.Alert, error)               { return m.FindOpenFn() }
func (m *AlertRepositoryMock) FindOpenByUser(userID uint) ([]entities.Alert, error) { return m.FindOpenByUserFn(userID) }
func (m *AlertRepositoryMock) Update(alert *entities.Alert) error                { return m.UpdateFn(alert) }
func (m *AlertRepositoryMock) Delete(id uint) error                              { return m.DeleteFn(id) }
func (m *AlertRepositoryMock) CountOpen() (int64, error)                         { return m.CountOpenFn() }
func (m *AlertRepositoryMock) CountOpenByUser(userID uint) (int64, error)        { return m.CountOpenByUserFn(userID) }
func (m *AlertRepositoryMock) CreateTx(tx ports.TxRunner, alert *entities.Alert) error { return m.CreateTxFn(tx, alert) }
func (m *AlertRepositoryMock) Create(alert *entities.Alert) error                { return m.CreateFn(alert) }
