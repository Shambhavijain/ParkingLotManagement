package inmemmory

import (
	"fmt"
	"parkingSlotManagement/internals/core/domain"
)

type SlotInMemmory struct {
	slots map[int]*domain.Slot
}

func NewSlotInMemmory() *SlotInMemmory {
	return &SlotInMemmory{
		slots: make(map[int]*domain.Slot)}
}

func (s *SlotInMemmory) SaveSlot(slot domain.Slot) error {
	s.slots[slot.SlotId] = &slot
	return nil
}
func (s *SlotInMemmory) UpdateSlot(slot *domain.Slot) error {
	existSlot, ok := s.slots[slot.SlotId]
	if !ok {
		return fmt.Errorf("slot of this id %d  not exists", slot.SlotId)
	}
	existSlot.IsFree = slot.IsFree
	existSlot.SlotType = slot.SlotType
	return nil
}
func (s *SlotInMemmory) ListAvailableSlots() ([]domain.Slot, error) {
	var availableSlots []domain.Slot
	for _, slot := range s.slots {
		if slot.IsFree {
			availableSlots = append(availableSlots, *slot)
		}

	}
	return availableSlots, nil
}
func (s *SlotInMemmory) FindSlotByType(SlotType string) ([]domain.Slot, error) {
	var availableSlots []domain.Slot
	for _, slot := range s.slots {
		if slot.SlotType == SlotType {
			availableSlots = append(availableSlots, *slot)
		}
	}
	return availableSlots, nil
}
func (s *SlotInMemmory) FindSlotTypebyID(SlotID int) string {

	existSlot, ok := s.slots[SlotID]
	if !ok {
		return "not exists slot of this slotid"
	}
	return existSlot.SlotType
}
func (s *SlotInMemmory) FindSlotByID(SlotId int) (*domain.Slot, error) {
	existsSlot, ok := s.slots[SlotId]
	if !ok {
		return nil, fmt.Errorf("slot of %d id not exists", SlotId)
	}
	return existsSlot, nil
}
