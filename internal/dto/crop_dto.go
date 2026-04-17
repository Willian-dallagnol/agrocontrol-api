package dto

type CreateCropRequest struct {
	Name    string `json:"name" binding:"required"`
	Type    string `json:"type"`
	FieldID uint   `json:"field_id" binding:"required"`
}

type UpdateCropRequest struct {
	Name    string `json:"name" binding:"required"`
	Type    string `json:"type"`
	FieldID uint   `json:"field_id" binding:"required"`
}

type CropResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	FieldID uint   `json:"field_id"`
}
