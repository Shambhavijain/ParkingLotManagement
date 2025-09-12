package ports

import "parkingSlotManagement/internals/core/domain"

type UserRepository interface {
	GetByUsername(username string) (*domain.Admin, error)
}
