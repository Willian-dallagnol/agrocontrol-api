package utils

import (
	"agrocontrol-api/internal/apperrors"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse é o contrato fixo de erro da API — nunca expõe internos
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// SuccessResponse envolve payloads de sucesso com campo "data"
type SuccessResponse struct {
	Data any `json:"data"`
}

func respondError(c *gin.Context, status int, message, code string) {
	c.JSON(status, ErrorResponse{Error: message, Code: code})
}

func RespondBadRequest(c *gin.Context, message string) {
	respondError(c, http.StatusBadRequest, message, "bad_request")
}

func RespondNotFound(c *gin.Context, message string) {
	respondError(c, http.StatusNotFound, message, "not_found")
}

func RespondUnauthorized(c *gin.Context, message string) {
	respondError(c, http.StatusUnauthorized, message, "unauthorized")
}

func RespondForbidden(c *gin.Context, message string) {
	respondError(c, http.StatusForbidden, message, "forbidden")
}

func RespondInternalError(c *gin.Context, message string) {
	respondError(c, http.StatusInternalServerError, message, "internal_error")
}

func RespondConflict(c *gin.Context, message string) {
	respondError(c, http.StatusConflict, message, "conflict")
}

func RespondUnprocessable(c *gin.Context, message string) {
	respondError(c, http.StatusUnprocessableEntity, message, "unprocessable")
}

// RespondDomainError traduz erros de domínio para HTTP de forma centralizada.
// Garante que erros internos nunca vazem para o cliente.
func RespondDomainError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		RespondNotFound(c, err.Error())
	case errors.Is(err, apperrors.ErrForbidden):
		RespondForbidden(c, err.Error())
	case errors.Is(err, apperrors.ErrUnauthorized):
		RespondUnauthorized(c, err.Error())
	case errors.Is(err, apperrors.ErrConflict):
		RespondConflict(c, err.Error())
	case errors.Is(err, apperrors.ErrInsufficientStock),
		errors.Is(err, apperrors.ErrInactiveField),
		errors.Is(err, apperrors.ErrNoActivePlanting),
		errors.Is(err, apperrors.ErrAlreadyHarvested),
		errors.Is(err, apperrors.ErrInvalidInput):
		RespondUnprocessable(c, err.Error())
	default:
		// Nunca expõe o erro interno ao cliente — apenas loga
		RespondInternalError(c, "erro interno do servidor")
	}
}
