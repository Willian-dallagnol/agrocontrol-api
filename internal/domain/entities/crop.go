package entities

import "gorm.io/gorm"

type Crop struct {
	gorm.Model
	Name    string
	Type    string
	FieldID uint
}
