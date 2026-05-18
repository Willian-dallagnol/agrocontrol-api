package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"agrocontrol-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// Login godoc
// @Summary      Autenticar usuário
// @Description  Autentica com email e senha, retorna access token (curta duração) e refresh token (7 dias).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Credenciais"
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} utils.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	resp, err := h.service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "code": "invalid_credentials"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// RefreshToken godoc
// @Summary      Renovar access token
// @Description  Recebe um refresh token válido e retorna um novo access token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshTokenRequest true "Refresh token"
// @Success      200 {object} dto.RefreshTokenResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} utils.ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	resp, err := h.service.RefreshToken(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "code": "invalid_refresh_token"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
