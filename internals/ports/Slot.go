package ports

import "parkingSlotManagement/internals/core/domain"

type SlotRepository interface {
	SaveSlot(slot domain.Slot) error
	UpdateSlot(slot *domain.Slot) error
	ListAvailableSlots() ([]domain.Slot, error)
	FindSlotByType(slottype string) ([]domain.Slot, error)
	FindSlotTypebyID(SlotId int) string
	FindSlotByID(SlotId int) (*domain.Slot, error)
}
