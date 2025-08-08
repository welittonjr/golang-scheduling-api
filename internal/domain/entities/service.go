package entities

import (
	"errors"
	"time"
)

type Service struct {
	id              int
	staffID         int
	name            string
	durationMinutes int
	price           float64
	createdAt       time.Time
}

func NewService(id, staffID int, name string, durationMinutes int, price float64) (*Service, error) {
	if name == "" {
		return nil, errors.New("nome do serviço é obrigatório")
	}
	if durationMinutes <= 0 {
		return nil, errors.New("a duração deve ser maior que zero")
	}
	if price < 0 {
		return nil, errors.New("preço não pode ser negativo")
	}

	return &Service{
		id:              id,
		staffID:         staffID,
		name:            name,
		durationMinutes: durationMinutes,
		price:           price,
		createdAt:       time.Now(),
	}, nil
}

func (s *Service) ID() int              { return s.id }
func (s *Service) StaffID() int         { return s.staffID }
func (s *Service) Name() string         { return s.name }
func (s *Service) DurationMinutes() int { return s.durationMinutes }
func (s *Service) Price() float64       { return s.price }
func (s *Service) CreatedAt() time.Time { return s.createdAt }

func (s *Service) ChangePrice(newPrice float64) error {
	if newPrice < 0 {
		return errors.New("preço não pode ser negativo")
	}
	s.price = newPrice
	return nil
}

func (s *Service) ChangeDuration(newDuration int) error {
	if newDuration <= 0 {
		return errors.New("duração inválida")
	}
	s.durationMinutes = newDuration
	return nil
}
