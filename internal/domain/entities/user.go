package entities

import "time"

// Role representa os papéis de usuário no sistema
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleManager  Role = "manager"
	RoleOperator Role = "operator"
)

// ValidRoles lista todas as roles aceitas pelo sistema
var ValidRoles = map[Role]struct{}{
	RoleAdmin: {}, RoleManager: {}, RoleOperator: {},
}

// User representa um usuário autenticado do sistema
type User struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         Role      `gorm:"not null;default:'operator'"`
	Active       bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// IsAdmin retorna true se o usuário tem role admin
func (u *User) IsAdmin() bool { return u.Role == RoleAdmin }

// IsManagerOrAbove retorna true se pode gerenciar recursos
func (u *User) IsManagerOrAbove() bool {
	return u.Role == RoleAdmin || u.Role == RoleManager
}
