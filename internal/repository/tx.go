package repository

import (
	"agrocontrol-api/internal/domain/ports"

	"gorm.io/gorm"
)

// GormTxRunner implementa ports.TxRunner usando *gorm.DB.
// É o único lugar do projeto que acopla transação ao GORM.
type GormTxRunner struct {
	db *gorm.DB
}

func NewGormTxRunner(db *gorm.DB) *GormTxRunner {
	return &GormTxRunner{db: db}
}

// DB expõe o *gorm.DB interno para uso pelos repositórios dentro da transação.
func (r *GormTxRunner) DB() *gorm.DB {
	return r.db
}

// RunInTx executa fn dentro de uma transação GORM.
// Se fn retornar erro, faz rollback. Caso contrário, commit.
func (r *GormTxRunner) RunInTx(fn func(tx ports.TxRunner) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(&GormTxRunner{db: tx})
	})
}
