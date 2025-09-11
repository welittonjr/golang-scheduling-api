package mocks

import (
	"time"

	"scheduling/internal/domain/entities"
)

type MockAvailableSlotRepository struct {
	FindByIDFunc                   func(id int) (*entities.AvailableSlot, error)
	FindAllByStaffIDFunc           func(staffID int) ([]*entities.AvailableSlot, error)
	FindSlotsByStaffAndDateFunc    func(staffID int, date time.Time) ([]*entities.AvailableSlot, error)
	HasConflictFunc                func(staffID int, weekday entities.Weekday, start, end time.Time) (bool, error)
	IsWithinAvailableSlotFunc      func(staffID int, start, end time.Time) (bool, error)
	SaveFunc                       func(slot *entities.AvailableSlot) error
	UpdateFunc                     func(slot *entities.AvailableSlot) error
	DeleteFunc                     func(id int) error
}

func (m *MockAvailableSlotRepository) FindByID(id int) (*entities.AvailableSlot, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockAvailableSlotRepository) FindAllByStaffID(staffID int) ([]*entities.AvailableSlot, error) {
	if m.FindAllByStaffIDFunc != nil {
		return m.FindAllByStaffIDFunc(staffID)
	}
	return nil, nil
}

func (m *MockAvailableSlotRepository) FindSlotsByStaffAndDate(staffID int, date time.Time) ([]*entities.AvailableSlot, error) {
	if m.FindSlotsByStaffAndDateFunc != nil {
		return m.FindSlotsByStaffAndDateFunc(staffID, date)
	}
	return nil, nil
}

func (m *MockAvailableSlotRepository) HasConflict(staffID int, weekday entities.Weekday, start, end time.Time) (bool, error) {
	if m.HasConflictFunc != nil {
		return m.HasConflictFunc(staffID, weekday, start, end)
	}
	return false, nil
}

func (m *MockAvailableSlotRepository) IsWithinAvailableSlot(staffID int, start, end time.Time) (bool, error) {
	if m.IsWithinAvailableSlotFunc != nil {
		return m.IsWithinAvailableSlotFunc(staffID, start, end)
	}
	return false, nil
}

func (m *MockAvailableSlotRepository) Save(slot *entities.AvailableSlot) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(slot)
	}
	return nil
}

func (m *MockAvailableSlotRepository) Update(slot *entities.AvailableSlot) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(slot)
	}
	return nil
}

func (m *MockAvailableSlotRepository) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func NewMockAvailableSlotRepository() *MockAvailableSlotRepository {
	return &MockAvailableSlotRepository{}
}