package parking_test

import (
	"testing"
	"time"

	"parkingSlotManagement/internals/core/domain"
	"parkingSlotManagement/internals/core/services/parking"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSlotRepo struct {
	mock.Mock
}

func (m *MockSlotRepo) FindSlotByType(vehicleType string) ([]domain.Slot, error) {
	args := m.Called(vehicleType)
	return args.Get(0).([]domain.Slot), args.Error(1)
}

func (m *MockSlotRepo) UpdateSlot(slot *domain.Slot) error {
	args := m.Called(slot)
	return args.Error(0)
}

func (m *MockSlotRepo) SaveSlot(slot domain.Slot) error {
	args := m.Called(slot)
	return args.Error(0)
}

func (m *MockSlotRepo) ListAvailableSlots() ([]domain.Slot, error) {
	args := m.Called()
	return args.Get(0).([]domain.Slot), args.Error(1)
}

func (m *MockSlotRepo) FindSlotByID(id int) (*domain.Slot, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Slot), args.Error(1)
}

func (m *MockSlotRepo) FindSlotTypebyID(id int) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

type MockTicketRepo struct {
	mock.Mock
}

func (m *MockTicketRepo) FindTicketByVehicleNumber(vehicleNumber string) (*domain.Ticket, error) {
	args := m.Called(vehicleNumber)
	return args.Get(0).(*domain.Ticket), args.Error(1)
}

func (m *MockTicketRepo) SaveTicket(ticket domain.Ticket) error {
	args := m.Called(ticket)
	return args.Error(0)
}

func (m *MockTicketRepo) DeleteTicket(ticketID int64) error {
	args := m.Called(ticketID)
	return args.Error(0)
}

func TestAddSlot(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)

	slot := domain.Slot{SlotId: 1, SlotType: "car", IsFree: true}
	slotRepo.On("SaveSlot", slot).Return(nil)

	err := service.AddSlot(slot)
	assert.NoError(t, err)
}

func TestGetAvailableSlots(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)

	expected := []domain.Slot{{SlotId: 1, SlotType: "car", IsFree: true}}
	slotRepo.On("ListAvailableSlots").Return(expected, nil)

	slots, err := service.GetAvailableSlots()
	assert.NoError(t, err)
	assert.Equal(t, expected, slots)
}

func TestCalculateFee_Car(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)

	entry := time.Now().Add(-2 * time.Hour)
	exit := time.Now()
	slotRepo.On("FindSlotTypebyID", 1).Return("car", nil)

	fee, err := service.CalculateFee(1, entry, exit)
	assert.NoError(t, err)
	assert.InDelta(t, 120.0, fee, 0.1)
}

func TestCalculateFee_Bike(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)

	entry := time.Now().Add(-2 * time.Hour)
	exit := time.Now()
	slotRepo.On("FindSlotTypebyID", 1).Return("bike", nil)

	fee, err := service.CalculateFee(1, entry, exit)
	assert.NoError(t, err)
	assert.InDelta(t, 60.0, fee, 0.1)
}

func TestUnparkVehicle(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)

	entry := time.Now().Add(-1 * time.Hour)
	ticket := &domain.Ticket{TicketId: 101, VehicleNumber: "UP16AB1234", SlotId: 1, EntryTime: entry}
	slot := &domain.Slot{SlotId: 1, SlotType: "car", IsFree: false}

	ticketRepo.On("FindTicketByVehicleNumber", "UP16AB1234").Return(ticket, nil)
	slotRepo.On("FindSlotByID", 1).Return(slot, nil)
	slotRepo.On("FindSlotTypebyID", 1).Return("car", nil)
	slotRepo.On("UpdateSlot", mock.Anything).Return(nil)
	ticketRepo.On("DeleteTicket", ticket.TicketId).Return(nil)

	fee, err := service.UnparkVehicle("UP16AB1234")
	assert.NoError(t, err)
	assert.True(t, fee > 0)
}
func TestParkVehicle_AlreadyParked(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)

	vehicle := domain.Vehicle{VehicleNumber: "UP16AB1234", VehicleType: "car"}
	existingTicket := &domain.Ticket{TicketId: 101, VehicleNumber: "UP16AB1234"}

	ticketRepo.On("FindTicketByVehicleNumber", vehicle.VehicleNumber).Return(existingTicket, nil)

	ticket, err := service.ParkVehicle(vehicle)
	assert.Nil(t, ticket)
	assert.Equal(t, parking.ErrVehicleAlreadyParked, err)
}

func TestParkVehicle_ExistingTicketCheckError(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)

	vehicle := domain.Vehicle{VehicleNumber: "UP16AB1234", VehicleType: "car"}
	ticketRepo.On("FindTicketByVehicleNumber", vehicle.VehicleNumber).Return((*domain.Ticket)(nil), assert.AnError)

	ticket, err := service.ParkVehicle(vehicle)
	assert.Nil(t, ticket)
	assert.Equal(t, parking.ErrExistingTicketCheck, err)
}

func TestParkVehicle_NoAvailableSlot(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)

	vehicle := domain.Vehicle{VehicleNumber: "UP16AB1234", VehicleType: "car"}
	ticketRepo.On("FindTicketByVehicleNumber", vehicle.VehicleNumber).Return((*domain.Ticket)(nil), nil)

	slotRepo.On("FindSlotByType", vehicle.VehicleType).Return([]domain.Slot{{SlotId: 1, SlotType: "car", IsFree: false}}, nil)

	ticket, err := service.ParkVehicle(vehicle)
	assert.Nil(t, ticket)
	assert.Equal(t, parking.ErrSlotFetchByType, err)
}

func TestUnparkVehicle_TicketNotFound(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)

	ticketRepo.On("FindTicketByVehicleNumber", "UP16AB1234").Return((*domain.Ticket)(nil), nil)

	fee, err := service.UnparkVehicle("UP16AB1234")
	assert.Equal(t, 0.0, fee)
	assert.Equal(t, parking.ErrTicketNotFound, err)
}

func TestUnparkVehicle_SlotNotFound(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)

	ticket := &domain.Ticket{TicketId: 101, VehicleNumber: "UP16AB1234", SlotId: 1, EntryTime: time.Now().Add(-1 * time.Hour)}
	ticketRepo.On("FindTicketByVehicleNumber", "UP16AB1234").Return(ticket, nil)
	slotRepo.On("FindSlotByID", 1).Return((*domain.Slot)(nil), nil)

	fee, err := service.UnparkVehicle("UP16AB1234")
	assert.Equal(t, 0.0, fee)
	assert.Equal(t, parking.ErrSlotNotFound, err)
}
