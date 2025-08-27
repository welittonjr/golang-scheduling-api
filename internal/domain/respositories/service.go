package respositories

import "scheduling/internal/domain/entities"

type ServiceRepository interface {
	FindByID(id int) (*entities.Service, error)
	FindAllByStaffID(staffID int) ([]*entities.Service, error)
	Exists(id int) (bool, error)
}
