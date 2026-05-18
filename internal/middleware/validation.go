package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationErrorResponse padroniza erros de validação
type ValidationErrorResponse struct {
	Error  string            `json:"error"`
	Code   string            `json:"code"`
	Fields map[string]string `json:"fields,omitempty"`
}

// FormatValidationErrors converte erros do validator em mensagens legíveis em português.
// Centraliza a lógica de formatação — handlers não precisam mais tratar isso individualmente.
func FormatValidationErrors(err error) (int, ValidationErrorResponse) {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return http.StatusBadRequest, ValidationErrorResponse{
			Error: err.Error(),
			Code:  "bad_request",
		}
	}

	fields := make(map[string]string, len(ve))
	for _, fe := range ve {
		fields[toSnakeCase(fe.Field())] = fieldErrorMessage(fe)
	}

	return http.StatusBadRequest, ValidationErrorResponse{
		Error:  "dados inválidos — verifique os campos",
		Code:   "validation_error",
		Fields: fields,
	}
}

// fieldErrorMessage retorna mensagem de erro em português para cada tag de validação
func fieldErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "campo obrigatório"
	case "min":
		if fe.Type().Kind().String() == "string" {
			return fmt.Sprintf("mínimo %s caracteres", fe.Param())
		}
		return fmt.Sprintf("valor mínimo: %s", fe.Param())
	case "max":
		if fe.Type().Kind().String() == "string" {
			return fmt.Sprintf("máximo %s caracteres", fe.Param())
		}
		return fmt.Sprintf("valor máximo: %s", fe.Param())
	case "gt":
		return fmt.Sprintf("deve ser maior que %s", fe.Param())
	case "gte":
		return fmt.Sprintf("deve ser maior ou igual a %s", fe.Param())
	case "lt":
		return fmt.Sprintf("deve ser menor que %s", fe.Param())
	case "lte":
		return fmt.Sprintf("deve ser menor ou igual a %s", fe.Param())
	case "email":
		return "e-mail inválido"
	case "len":
		return fmt.Sprintf("deve ter exatamente %s caracteres", fe.Param())
	case "oneof":
		return fmt.Sprintf("valor inválido — opções: %s", fe.Param())
	case "url":
		return "URL inválida"
	case "numeric":
		return "deve ser numérico"
	case "alpha":
		return "deve conter apenas letras"
	case "alphanum":
		return "deve conter apenas letras e números"
	default:
		return fmt.Sprintf("validação falhou: %s", fe.Tag())
	}
}

// toSnakeCase converte PascalCase/camelCase para snake_case para os nomes dos campos
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
