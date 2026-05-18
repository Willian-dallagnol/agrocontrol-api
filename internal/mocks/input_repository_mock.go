package mocks

import (
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/domain/ports"
)

type InputRepositoryMock struct {
	CreateFn              func(input *entities.Input) error
	FindByIDFn            func(id uint) (*entities.Input, error)
	FindAllPaginatedFn    func(offset, limit int, search, category string) ([]entities.Input, int64, error)
	FindByUserPaginatedFn func(userID uint, offset, limit int, search, category string) ([]entities.Input, int64, error)
	FindAllFn             func() ([]entities.Input, error)
	FindByUserFn          func(userID uint) ([]entities.Input, error)
	FindLowStockFn        func() ([]entities.Input, error)
	FindExpiringSoonFn    func(days int) ([]entities.Input, error)
	UpdateFn              func(input *entities.Input) error
	DeleteFn              func(id uint) error
	CountLowStockFn       func() (int64, error)
	CreateTxFn            func(tx ports.TxRunner, input *entities.Input) error
	DeductStockTxFn       func(tx ports.TxRunner, id uint, qty float64) error
	FindByIDTxFn          func(tx ports.TxRunner, id uint) (*entities.Input, error)
}

var _ ports.InputRepository = (*InputRepositoryMock)(nil)

func (m *InputRepositoryMock) Create(i *entities.Input) error  { return m.CreateFn(i) }
func (m *InputRepositoryMock) FindByID(id uint) (*entities.Input, error) { return m.FindByIDFn(id) }
func (m *InputRepositoryMock) FindAllPaginated(o, l int, s, c string) ([]entities.Input, int64, error) { return m.FindAllPaginatedFn(o, l, s, c) }
func (m *InputRepositoryMock) FindByUserPaginated(uID uint, o, l int, s, c string) ([]entities.Input, int64, error) { return m.FindByUserPaginatedFn(uID, o, l, s, c) }
func (m *InputRepositoryMock) FindAll() ([]entities.Input, error) { return m.FindAllFn() }
func (m *InputRepositoryMock) FindByUser(userID uint) ([]entities.Input, error) { return m.FindByUserFn(userID) }
func (m *InputRepositoryMock) FindLowStock() ([]entities.Input, error) { return m.FindLowStockFn() }
func (m *InputRepositoryMock) FindExpiringSoon(days int) ([]entities.Input, error) { return m.FindExpiringSoonFn(days) }
func (m *InputRepositoryMock) Update(i *entities.Input) error  { return m.UpdateFn(i) }
func (m *InputRepositoryMock) Delete(id uint) error            { return m.DeleteFn(id) }
func (m *InputRepositoryMock) CountLowStock() (int64, error)   { return m.CountLowStockFn() }
func (m *InputRepositoryMock) CreateTx(tx ports.TxRunner, i *entities.Input) error { return m.CreateTxFn(tx, i) }
func (m *InputRepositoryMock) DeductStockTx(tx ports.TxRunner, id uint, qty float64) error { return m.DeductStockTxFn(tx, id, qty) }
func (m *InputRepositoryMock) FindByIDTx(tx ports.TxRunner, id uint) (*entities.Input, error) { return m.FindByIDTxFn(tx, id) }
