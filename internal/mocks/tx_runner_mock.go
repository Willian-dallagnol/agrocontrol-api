package mocks

import "agrocontrol-api/internal/domain/ports"

// TxRunnerMock executa a função diretamente sem transação real.
// Útil para testes de serviço que precisam de TxRunner.
type TxRunnerMock struct{}

var _ ports.TxRunner = (*TxRunnerMock)(nil)

func (m *TxRunnerMock) RunInTx(fn func(tx ports.TxRunner) error) error {
	return fn(m)
}
