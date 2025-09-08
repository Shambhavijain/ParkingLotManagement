package parking

import (
	"database/sql"
	"parkingSlotManagement/internals/core/domain"
	"parkingSlotManagement/internals/ports"
	"time"
)

type ParkingService struct {
	SlotRepo   ports.SlotRepository
	TicketRepo ports.TicketRepository
}

func NewParkingService(s ports.SlotRepository, t ports.TicketRepository) *ParkingService {
	return &ParkingService{SlotRepo: s,
		TicketRepo: t,
	}
}

func (s *ParkingService) ParkVehicle(vehicle domain.Vehicle) (*domain.Ticket, error) {

	existingTicket, err := s.TicketRepo.FindTicketByVehicleNumber(vehicle.VehicleNumber)
	if err != nil && err != sql.ErrNoRows {
		return nil, ErrExistingTicketCheck
	}
	if existingTicket != nil {

		return nil, ErrVehicleAlreadyParked
	}

	slots, err := s.SlotRepo.FindSlotByType(vehicle.VehicleType)
	if err != nil {
		return nil, ErrSlotListFailed
	}
	var firstAvailable *domain.Slot
	for i := range slots {
		if slots[i].IsFree {
			firstAvailable = &slots[i]
			break

		}
	}
	if firstAvailable == nil {
		return nil, ErrSlotFetchByType
	}
	firstAvailable.IsFree = false
	if err := s.SlotRepo.UpdateSlot(firstAvailable); err != nil {
		return nil, ErrSlotUpdateFailed
	}
	ticket := &domain.Ticket{
		TicketId:      GenerateTicketID(),
		VehicleNumber: vehicle.VehicleNumber,
		SlotId:        firstAvailable.SlotId,
		EntryTime:     time.Now(),
	}
	if err := s.TicketRepo.SaveTicket(*ticket); err != nil {
		return nil, ErrTicketSaveFailed
	}
	return ticket, nil

}

func GenerateTicketID() int64 {
	return time.Now().UnixNano()

}

func (s *ParkingService) UnparkVehicle(VehicleNumber string) (float64, error) {
	ExitTime := time.Now()
	ticket, err := s.TicketRepo.FindTicketByVehicleNumber(VehicleNumber)
	if err != nil || ticket == nil {
		return 0, ErrTicketNotFound
	}
	slot, err := s.SlotRepo.FindSlotByID(ticket.SlotId)

	if err != nil || slot == nil {
		return 0, ErrSlotNotFound
	}

	fee, err := s.CalculateFee(ticket.SlotId, ticket.EntryTime, ExitTime)
	if err != nil {
		return 0, ErrFeeCalculationFailed
	}

	slot.IsFree = true
	if err := s.SlotRepo.UpdateSlot(slot); err != nil {
		return 0, ErrSlotUpdateFailed
	}

	if err = s.TicketRepo.DeleteTicket(ticket.TicketId); err != nil {
		return 0, ErrTicketDeleteFailed
	}

	return fee, nil

}
func (s *ParkingService) AddSlot(slot domain.Slot) error {
	err := s.SlotRepo.SaveSlot(slot)
	return err

}
func (s *ParkingService) GetAvailableSlots() ([]domain.Slot, error) {
	slots, err := s.SlotRepo.ListAvailableSlots()
	if err != nil {
		return nil, ErrSlotListFailed
	}
	return slots, nil
}

func (s *ParkingService) CalculateFee(SlotId int, EntryTime time.Time, ExistTime time.Time) (float64, error) {
	slottype, err := s.SlotRepo.FindSlotTypebyID(SlotId)
	if err != nil {
		return 0, err
	}
	duration := ExistTime.Sub(EntryTime)

	switch slottype {
	case "car":
		fee := duration.Hours() * 60
		return fee, nil
	case "bike":
		fee := duration.Hours() * 30
		return fee, nil
	default:
		return 0, ErrInvalidVehicleType
	}
}
