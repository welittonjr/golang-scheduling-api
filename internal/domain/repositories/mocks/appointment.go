package mocks

import (
	"time"

	"scheduling/internal/domain/entities"
)

type MockAppointmentRepository struct {
	FindByIDFunc           func(id int) (*entities.Appointment, error)
	FindAllByStaffIDFunc   func(staffID int) ([]*entities.Appointment, error)
	HasConflictFunc        func(staffID int, start, end time.Time) (bool, error)
	SaveFunc               func(appointment *entities.Appointment) error
	UpdateFunc             func(appointment *entities.Appointment) error
	DeleteFunc             func(id int) error
}

func (m *MockAppointmentRepository) FindByID(id int) (*entities.Appointment, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockAppointmentRepository) FindAllByStaffID(staffID int) ([]*entities.Appointment, error) {
	if m.FindAllByStaffIDFunc != nil {
		return m.FindAllByStaffIDFunc(staffID)
	}
	return nil, nil
}

func (m *MockAppointmentRepository) HasConflict(staffID int, start, end time.Time) (bool, error) {
	if m.HasConflictFunc != nil {
		return m.HasConflictFunc(staffID, start, end)
	}
	return false, nil
}

func (m *MockAppointmentRepository) Save(appointment *entities.Appointment) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(appointment)
	}
	return nil
}

func (m *MockAppointmentRepository) Update(appointment *entities.Appointment) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(appointment)
	}
	return nil
}

func (m *MockAppointmentRepository) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func NewMockAppointmentRepository() *MockAppointmentRepository {
	return &MockAppointmentRepository{}
}