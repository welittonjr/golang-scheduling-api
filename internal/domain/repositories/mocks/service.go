package mocks

import "scheduling/internal/domain/entities"

type MockServiceRepository struct {
	FindByIDFunc         func(id int) (*entities.Service, error)
	FindAllByStaffIDFunc func(staffID int) ([]*entities.Service, error)
	ExistsFunc           func(id int) (bool, error)
}

func (m *MockServiceRepository) FindByID(id int) (*entities.Service, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockServiceRepository) FindAllByStaffID(staffID int) ([]*entities.Service, error) {
	if m.FindAllByStaffIDFunc != nil {
		return m.FindAllByStaffIDFunc(staffID)
	}
	return nil, nil
}

func (m *MockServiceRepository) Exists(id int) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(id)
	}
	return false, nil
}

func NewMockServiceRepository() *MockServiceRepository {
	return &MockServiceRepository{}
}