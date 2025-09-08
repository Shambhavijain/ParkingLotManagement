package parking

import (
	"database/sql"

	"parkingSlotManagement/internals/adapters/repositories/inmemmory"
	"parkingSlotManagement/internals/core/domain"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParkVehicle(t *testing.T) {
	slotrepo := inmemmory.NewSlotInMemmory()
	ticketrepo := inmemmory.NewTicketInMemmory()
	slot := domain.Slot{
		SlotId:   1,
		SlotType: "car",
		IsFree:   true,
	}
	err := slotrepo.SaveSlot(slot)
	assert.NoError(t, err)
	service := NewParkingService(slotrepo, ticketrepo)
	vehicle := domain.Vehicle{
		VehicleNumber: "UP74M8311",
		VehicleType:   "car",
	}

	ticket, err := service.ParkVehicle(vehicle)

	err1 := ticketrepo.SaveTicket(*ticket)
	assert.NoError(t, err1)

	assert.NoError(t, err)
	assert.NotNil(t, ticket)

	assert.Equal(t, vehicle.VehicleNumber, ticket.VehicleNumber)
	assert.Equal(t, slot.SlotId, ticket.SlotId)

	ticket2, err2 := service.ParkVehicle(vehicle)
	assert.Error(t, err2)
	assert.Nil(t, ticket2)
	assert.Contains(t, err2.Error(), "already parked")

}

func TestGenerateTicketId(t *testing.T) {
	id1 := GenerateTicketID()
	time.Sleep(1 * time.Second)
	id2 := GenerateTicketID()

	assert.NotZero(t, id1)
	assert.NotZero(t, id2)
	assert.NotEqual(t, id1, id2)
}
func TestUnparkVehicle(t *testing.T) {
	slotRepo := inmemmory.NewSlotInMemmory()
	ticketRepo := inmemmory.NewTicketInMemmory()
	slot := domain.Slot{
		SlotId:   1,
		SlotType: "car",
		IsFree:   true,
	}

	_ = slotRepo.SaveSlot(slot)

	entryTime := time.Now().Add(-2 * time.Hour)

	ticket := domain.Ticket{
		TicketId:      123456987654321,
		VehicleNumber: "UP74M8311",
		SlotId:        1,
		EntryTime:     entryTime,
	}

	_ = ticketRepo.SaveTicket(ticket)

	service := NewParkingService(slotRepo, ticketRepo)
	fee, err := service.UnparkVehicle("UP74M8311")

	assert.NoError(t, err)
	assert.Greater(t, fee, float64(0))

	updatedSlot, _ := slotRepo.FindSlotByID(1)
	assert.True(t, updatedSlot.IsFree)

	_, err = ticketRepo.FindTicketByVehicleNumber("UP74M8311")
	assert.ErrorIs(t, err, sql.ErrNoRows)

	slot1 := domain.Slot{
		SlotId:   2,
		SlotType: "bike",
		IsFree:   true,
	}
	_ = slotRepo.SaveSlot(slot1)
	ticket1 := domain.Ticket{
		TicketId:      123456987654321,
		VehicleNumber: "UP74M8412",
		SlotId:        2,
		EntryTime:     entryTime,
	}
	_ = ticketRepo.SaveTicket(ticket1)
	fee1, err := service.UnparkVehicle("UP74M8412")
	assert.NoError(t, err)
	assert.Greater(t, fee1, float64(0))

	updatedSlot1, _ := slotRepo.FindSlotByID(2)
	assert.True(t, updatedSlot1.IsFree)

}
func TestAddSlot(t *testing.T) {
	slotRepo := inmemmory.NewSlotInMemmory()
	service := NewParkingService(slotRepo, nil)
	slot := domain.Slot{
		SlotId:   1,
		SlotType: "car",
		IsFree:   true,
	}
	err := service.AddSlot(slot)
	assert.NoError(t, err)

}
func TestGetAvailableSlots(t *testing.T) {
	slotRepo := inmemmory.NewSlotInMemmory()
	service := NewParkingService(slotRepo, nil)
	slots := []domain.Slot{
		{SlotId: 1, SlotType: "car", IsFree: true},
		{SlotId: 2, SlotType: "bus", IsFree: true},
		{SlotId: 3, SlotType: "bus", IsFree: false},
	}
	for _, slot := range slots {

		slotRepo.SaveSlot(slot)
	}
	availableSlots, err := service.GetAvailableSlots()
	assert.NoError(t, err)
	assert.Len(t, availableSlots, 2)
	for _, slot := range availableSlots {
		assert.True(t, slot.IsFree)
	}

}

func TestOpposingCases(t *testing.T) {
	tests := []struct {
		name        string
		ticket      *domain.Ticket
		slot        *domain.Slot
		expectedFee float64
		expectError bool
		errorText   string
	}{
		{
			name:        "ticket not found",
			ticket:      nil,
			slot:        nil,
			expectError: true,
			errorText:   "ticket of this vehiclenumber not found",
		},
		{
			name: "slot not found",
			ticket: &domain.Ticket{
				TicketId:      1,
				VehicleNumber: "XYZ123",
				SlotId:        101,
				EntryTime:     time.Now().Add(-2 * time.Hour),
			},
			slot:        nil,
			expectError: true,
			errorText:   "failed to fetch slot",
		},
		{
			name: "invalid slot type",
			ticket: &domain.Ticket{
				TicketId:      1,
				VehicleNumber: "XYZ123",
				SlotId:        101,
				EntryTime:     time.Now().Add(-2 * time.Hour),
			},
			slot: &domain.Slot{
				SlotId:   101,
				IsFree:   false,
				SlotType: "invalid",
			},
			expectError: true,
			errorText:   "unable to calculate fee",
		},
		{
			name: "success case",
			ticket: &domain.Ticket{
				TicketId:      1,
				VehicleNumber: "XYZ123",
				SlotId:        101,
				EntryTime:     time.Now().Add(-2 * time.Hour),
			},
			slot: &domain.Slot{
				SlotId:   101,
				IsFree:   false,
				SlotType: "car",
			},

			expectedFee: 120.0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticketRepo := inmemmory.NewTicketInMemmory()
			slotRepo := inmemmory.NewSlotInMemmory()
			service := NewParkingService(slotRepo, ticketRepo)

			if tt.ticket != nil {
				ticketRepo.SaveTicket(*tt.ticket)
			}

			if tt.slot != nil {
				slotRepo.SaveSlot(*tt.slot)
			}

			fee, err := service.UnparkVehicle("XYZ123")

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorText)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expectedFee, fee, 0.1)
			}
		})
	}
}
