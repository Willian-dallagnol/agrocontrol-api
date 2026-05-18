package apperrors

import (
	"errors"
	"fmt"
)

// ── Erros sentinela ───────────────────────────────────────────────────────────

var (
	ErrNotFound          = errors.New("registro não encontrado")
	ErrForbidden         = errors.New("acesso negado")
	ErrConflict          = errors.New("registro já existe")
	ErrInvalidInput      = errors.New("dados inválidos")
	ErrInsufficientStock = errors.New("estoque insuficiente")
	ErrInactiveField     = errors.New("talhão inativo")
	ErrNoActivePlanting  = errors.New("não há plantio ativo para este talhão")
	ErrAlreadyHarvested  = errors.New("este plantio já foi colhido")
	ErrUnauthorized      = errors.New("não autorizado")
)

// ── Construtores com contexto ─────────────────────────────────────────────────

func NotFoundError(resource string, id any) error {
	return fmt.Errorf("%s de id %v não encontrado: %w", resource, id, ErrNotFound)
}

func ForbiddenError(action string) error {
	return fmt.Errorf("não autorizado a %s: %w", action, ErrForbidden)
}

func ConflictError(field string) error {
	return fmt.Errorf("%s já cadastrado: %w", field, ErrConflict)
}

// ── Helpers de verificação ────────────────────────────────────────────────────

func IsNotFound(err error) bool         { return errors.Is(err, ErrNotFound) }
func IsForbidden(err error) bool        { return errors.Is(err, ErrForbidden) }
func IsConflict(err error) bool         { return errors.Is(err, ErrConflict) }
func IsInvalidInput(err error) bool     { return errors.Is(err, ErrInvalidInput) }
func IsInsufficientStock(err error) bool { return errors.Is(err, ErrInsufficientStock) }
