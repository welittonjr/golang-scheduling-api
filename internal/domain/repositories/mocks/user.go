package mocks

import "scheduling/internal/domain/entities"

type MockUserRepository struct {
	FindByIDFunc func(id int) (*entities.User, error)
	ExistsFunc   func(id int) (bool, error)
}

func (m *MockUserRepository) FindByID(id int) (*entities.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockUserRepository) Exists(id int) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(id)
	}
	return false, nil
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}