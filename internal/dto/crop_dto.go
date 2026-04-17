package dto

// 🌾 Estrutura usada para criar uma nova cultura (Crop)
type CreateCropRequest struct {
	Name string `json:"name" binding:"required"`
	// 👉 nome da cultura (ex: Soja, Milho)
	// 🔥 obrigatório

	Type string `json:"type"`
	// 👉 tipo da cultura (ex: grão, leguminosa)
	// (campo opcional)

	FieldID uint `json:"field_id" binding:"required"`
	// 👉 ID do talhão (Field) ao qual a cultura pertence
	// 🔗 obrigatório para manter o relacionamento
}

// 🔄 Estrutura usada para atualizar uma cultura existente
type UpdateCropRequest struct {
	Name string `json:"name" binding:"required"`
	// 👉 novo nome da cultura

	Type string `json:"type"`
	// 👉 novo tipo da cultura

	FieldID uint `json:"field_id" binding:"required"`
	// 👉 pode atualizar o vínculo com outro talhão
}

// 📤 Estrutura de resposta da API (o que o usuário recebe)
type CropResponse struct {
	ID uint `json:"id"`
	// 👉 identificador da cultura

	Name string `json:"name"`
	// 👉 nome da cultura

	Type string `json:"type"`
	// 👉 tipo da cultura

	FieldID uint `json:"field_id"`
	// 👉 referência ao talhão relacionado
}
