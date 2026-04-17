package dto

// 🌱 Estrutura usada para criar um novo talhão (Field)
type CreateFieldRequest struct {
	Name string `json:"name" binding:"required"`
	// 👉 nome do talhão (ex: Talhão 1)
	// 🔥 obrigatório

	Area float64 `json:"area"`
	// 👉 área do talhão
	// ⚠️ validação (ex: > 0) é feita no Service

	SoilType string `json:"soil_type"`
	// 👉 tipo de solo (argiloso, arenoso, etc)
	// (campo opcional)

	FarmID uint `json:"farm_id" binding:"required"`
	// 👉 ID da fazenda à qual o talhão pertence
	// 🔗 obrigatório para manter o relacionamento
}

// 🔄 Estrutura usada para atualizar um talhão existente
type UpdateFieldRequest struct {
	Name string `json:"name" binding:"required"`
	// 👉 novo nome do talhão

	Area float64 `json:"area"`
	// 👉 nova área (validada no Service)

	SoilType string `json:"soil_type"`
	// 👉 novo tipo de solo

	FarmID uint `json:"farm_id" binding:"required"`
	// 👉 permite alterar o vínculo com outra fazenda
}

// 📤 Estrutura de resposta da API
type FieldResponse struct {
	ID uint `json:"id"`
	// 👉 identificador do talhão

	Name string `json:"name"`
	// 👉 nome do talhão

	Area float64 `json:"area"`
	// 👉 área do talhão

	SoilType string `json:"soil_type"`
	// 👉 tipo de solo

	FarmID uint `json:"farm_id"`
	// 👉 referência à fazenda
}
