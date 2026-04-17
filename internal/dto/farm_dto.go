package dto

// 🚜 Estrutura usada para criar uma nova fazenda (entrada da API)
type CreateFarmRequest struct {
	Name string `json:"name" binding:"required"`
	// 👉 nome da fazenda (obrigatório)

	OwnerName string `json:"owner_name" binding:"required"`
	// 👉 nome do proprietário (obrigatório)

	Location string `json:"location"`
	// 👉 localização geral (opcional)

	TotalArea float64 `json:"total_area"`
	// 👉 área total da fazenda
	// ⚠️ validação de valor (> 0) é feita no Service

	City string `json:"city" binding:"required"`
	// 👉 cidade (obrigatório)

	State string `json:"state" binding:"required"`
	// 👉 estado (obrigatório)
}

// 🔄 Estrutura usada para atualizar uma fazenda existente
type UpdateFarmRequest struct {
	Name string `json:"name" binding:"required"`
	// 👉 novo nome da fazenda

	OwnerName string `json:"owner_name" binding:"required"`
	// 👉 novo proprietário

	Location string `json:"location"`
	// 👉 nova localização

	TotalArea float64 `json:"total_area"`
	// 👉 nova área total (validada no Service)

	City string `json:"city" binding:"required"`
	// 👉 cidade

	State string `json:"state" binding:"required"`
	// 👉 estado
}

// 📤 Estrutura de resposta da API
type FarmResponse struct {
	ID uint `json:"id"`
	// 👉 identificador da fazenda

	Name string `json:"name"`
	// 👉 nome da fazenda

	OwnerName string `json:"owner_name"`
	// 👉 proprietário

	Location string `json:"location"`
	// 👉 localização

	TotalArea float64 `json:"total_area"`
	// 👉 área total

	City string `json:"city"`
	// 👉 cidade

	State string `json:"state"`
	// 👉 estado
}
