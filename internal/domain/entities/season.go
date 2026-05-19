package entities

import "time"

type SeasonStatus string

const (
	SeasonStatusPlanning SeasonStatus = "planning"
	SeasonStatusActive   SeasonStatus = "active"
	SeasonStatusFinished SeasonStatus = "finished"
)

type Season struct {
	ID        uint         `gorm:"primaryKey"`
	Name      string       `gorm:"not null;index"`
	StartDate time.Time    `gorm:"not null"`
	EndDate   time.Time    `gorm:"not null"`
	Status    SeasonStatus `gorm:"not null;default:'planning';index"`
	CreatedBy uint         `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsActive retorna true se a safra está em andamento.
func (s *Season) IsActive() bool {
	return s.Status == SeasonStatusActive
}

// IsFinished retorna true se a safra foi encerrada.
func (s *Season) IsFinished() bool {
	return s.Status == SeasonStatusFinished
}

// DurationDays retorna a duração da safra em dias.
func (s *Season) DurationDays() int {
	return int(s.EndDate.Sub(s.StartDate).Hours() / 24)
}

// IsOngoing verifica se a data atual está dentro do período da safra.
func (s *Season) IsOngoing() bool {
	now := time.Now()
	return now.After(s.StartDate) && now.Before(s.EndDate)
}

// Activate muda o status da safra para ativa.
func (s *Season) Activate() {
	s.Status = SeasonStatusActive
}

// Finish encerra a safra.
func (s *Season) Finish() {
	s.Status = SeasonStatusFinished
}
