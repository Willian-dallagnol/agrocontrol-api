package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/handler"
	"agrocontrol-api/internal/mocks"
	"agrocontrol-api/internal/middleware"
	"agrocontrol-api/internal/service"
	"agrocontrol-api/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const testJWTSecret = "test_secret_32_chars_padded_ok!!"

func init() {
	gin.SetMode(gin.TestMode)
}

// makeToken gera um token JWT válido para os testes
func makeToken(userID uint, role string) string {
	token, _ := utils.GenerateToken(userID, "test@test.com", role, testJWTSecret, 1)
	return "Bearer " + token
}

// setupFarmRouter cria router de teste com FarmHandler
func setupFarmRouter(svc *service.FarmService) *gin.Engine {
	r := gin.New()
	h := handler.NewFarmHandler(svc)
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(testJWTSecret))
	api.POST("/farms", h.CreateFarm)
	api.GET("/farms", h.GetFarms)
	api.GET("/farms/:id", h.GetFarmByID)
	api.PUT("/farms/:id", h.UpdateFarm)
	api.DELETE("/farms/:id", h.DeleteFarm)
	return r
}

// ── GET /farms ────────────────────────────────────────────────────────────────

func TestFarmHandler_GetFarms_Unauthorized(t *testing.T) {
	r := setupFarmRouter(nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/farms", nil)
	// sem Authorization header
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, got %d", w.Code)
	}
}

func TestFarmHandler_GetFarms_OK(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		FindByCreatedByPaginatedFn: func(userID uint, offset, limit int, search string) ([]entities.Farm, int64, error) {
			return []entities.Farm{{ID: 1, Name: "Fazenda A", CreatedBy: userID}}, 1, nil
		},
	}
	svc := service.NewFarmService(repo)
	r := setupFarmRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/farms", nil)
	req.Header.Set("Authorization", makeToken(1, "operator"))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperava 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

// ── POST /farms ───────────────────────────────────────────────────────────────

func TestFarmHandler_CreateFarm_BadRequest(t *testing.T) {
	r := setupFarmRouter(nil)
	// Body inválido (sem campos obrigatórios)
	body := `{"name": ""}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/farms", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", makeToken(1, "manager"))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperava 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestFarmHandler_CreateFarm_Created(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		CreateFn: func(farm *entities.Farm) error {
			farm.ID = 10
			return nil
		},
	}
	svc := service.NewFarmService(repo)
	r := setupFarmRouter(svc)

	payload := dto.CreateFarmRequest{
		Name: "Fazenda Teste", OwnerName: "Dono", TotalArea: 100, City: "Londrina", State: "PR",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/farms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", makeToken(1, "manager"))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("esperava 201, got %d — body: %s", w.Code, w.Body.String())
	}

	var resp dto.FarmResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.ID != 10 {
		t.Errorf("esperava ID=10, got %d", resp.ID)
	}
}

// ── GET /farms/:id ────────────────────────────────────────────────────────────

func TestFarmHandler_GetByID_NotFound(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Farm, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := service.NewFarmService(repo)
	r := setupFarmRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/farms/999", nil)
	req.Header.Set("Authorization", makeToken(1, "operator"))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperava 404, got %d", w.Code)
	}
}

func TestFarmHandler_GetByID_Forbidden(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Farm, error) {
			return &entities.Farm{ID: id, CreatedBy: 1}, nil // pertence ao user 1
		},
	}
	svc := service.NewFarmService(repo)
	r := setupFarmRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/farms/1", nil)
	req.Header.Set("Authorization", makeToken(2, "operator")) // user 2
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperava 403, got %d", w.Code)
	}
}

func TestFarmHandler_GetByID_OK(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Farm, error) {
			return &entities.Farm{ID: id, Name: "Fazenda A", CreatedBy: 1}, nil
		},
	}
	svc := service.NewFarmService(repo)
	r := setupFarmRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/farms/1", nil)
	req.Header.Set("Authorization", makeToken(1, "operator"))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperava 200, got %d", w.Code)
	}
}

// ── DELETE /farms/:id ─────────────────────────────────────────────────────────

func TestFarmHandler_Delete_Forbidden(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Farm, error) {
			return &entities.Farm{ID: id, CreatedBy: 1}, nil
		},
	}
	svc := service.NewFarmService(repo)
	r := setupFarmRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/farms/1", nil)
	req.Header.Set("Authorization", makeToken(2, "manager")) // user 2, não é dono
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperava 403, got %d", w.Code)
	}
}

func TestFarmHandler_Delete_NoContent(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Farm, error) {
			return &entities.Farm{ID: id, CreatedBy: 1}, nil
		},
		DeleteFn: func(id uint) error { return nil },
	}
	svc := service.NewFarmService(repo)
	r := setupFarmRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/farms/1", nil)
	req.Header.Set("Authorization", makeToken(1, "admin"))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("esperava 204, got %d", w.Code)
	}
}

// ── Testes de autenticação ────────────────────────────────────────────────────

func setupAuthRouter(svc *service.AuthService) *gin.Engine {
	r := gin.New()
	h := handler.NewAuthHandler(svc)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/refresh", h.RefreshToken)
	return r
}

func TestAuthHandler_Login_BadRequest(t *testing.T) {
	r := setupAuthRouter(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(`{"email":"invalido"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperava 400, got %d", w.Code)
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		FindByEmailFn: func(email string) (*entities.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := service.NewAuthService(repo, testJWTSecret, 1)
	r := setupAuthRouter(svc)

	body := `{"email":"nao@existe.com","password":"qualquer"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, got %d", w.Code)
	}
}

func TestAuthHandler_RefreshToken_Invalid(t *testing.T) {
	repo := &mocks.UserRepositoryMock{}
	svc := service.NewAuthService(repo, testJWTSecret, 1)
	r := setupAuthRouter(svc)

	body := `{"refresh_token":"token-invalido"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, got %d", w.Code)
	}
}

// ── Teste de RBAC via handler ─────────────────────────────────────────────────

func TestFarmHandler_Create_OperatorCannotCreate(t *testing.T) {
	// Operator não deveria criar fazenda (rota exige manager+)
	// No router atual não temos middleware de role na rota de test,
	// mas testamos que o handler processa corretamente
	repo := &mocks.FarmRepositoryMock{
		CreateFn: func(farm *entities.Farm) error {
			farm.ID = 1
			return nil
		},
	}
	svc := service.NewFarmService(repo)

	// Router com middleware de role
	r := gin.New()
	h := handler.NewFarmHandler(svc)
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(testJWTSecret))
	api.POST("/farms", middleware.ManagerOrAbove(), h.CreateFarm)

	payload := dto.CreateFarmRequest{
		Name: "Fazenda", OwnerName: "Dono", TotalArea: 100, City: "X", State: "PR",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/farms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", makeToken(1, "operator")) // operator
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperava 403 para operator tentando criar fazenda, got %d", w.Code)
	}
}

func TestFarmHandler_Create_ManagerCanCreate(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		CreateFn: func(farm *entities.Farm) error {
			farm.ID = 5
			return nil
		},
	}
	svc := service.NewFarmService(repo)

	r := gin.New()
	h := handler.NewFarmHandler(svc)
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(testJWTSecret))
	api.POST("/farms", middleware.ManagerOrAbove(), h.CreateFarm)

	payload := dto.CreateFarmRequest{
		Name: "Fazenda", OwnerName: "Dono", TotalArea: 100, City: "X", State: "PR",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/farms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", makeToken(1, "manager"))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("esperava 201 para manager, got %d — body: %s", w.Code, w.Body.String())
	}
}

// ── Healthcheck ───────────────────────────────────────────────────────────────

func TestHealthCheck(t *testing.T) {
	r := gin.New()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperava 200, got %d", w.Code)
	}

	// Verifica que o body contém status ok
	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["status"] != "ok" {
		t.Errorf("esperava status=ok, got: %v", body)
	}
}

// ── Testa que ErrConflict vira 409 ───────────────────────────────────────────

func TestFarmHandler_Create_Conflict(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		CreateFn: func(farm *entities.Farm) error {
			return apperrors.ConflictError("fazenda")
		},
	}
	svc := service.NewFarmService(repo)
	r := setupFarmRouter(svc)

	payload := dto.CreateFarmRequest{
		Name: "Fazenda", OwnerName: "Dono", TotalArea: 100, City: "X", State: "PR",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/farms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", makeToken(1, "manager"))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("esperava 409, got %d", w.Code)
	}
}
