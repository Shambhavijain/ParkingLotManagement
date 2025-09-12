package requestHandlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestParkVehicleRequest(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)
	handler := NewHandlers(service)

	vehicle := domain.Vehicle{VehicleNumber: "UP16AB1234", VehicleType: "car"}
	slots := []domain.Slot{{SlotId: 1, SlotType: "car", IsFree: true}}
	ticketRepo.On("FindTicketByVehicleNumber", vehicle.VehicleNumber).Return((*domain.Ticket)(nil), nil)

	slotRepo.On("FindSlotByType", vehicle.VehicleType).Return(slots, nil)
	slotRepo.On("UpdateSlot", &slots[0]).Return(nil)
	ticketRepo.On("SaveTicket", mock.AnythingOfType("domain.Ticket")).Return(nil)

	body, _ := json.Marshal(vehicle)
	req := httptest.NewRequest(http.MethodPost, "/park", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ParkVehicleRequest(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestUnparkVehicleRequest(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)
	handler := NewHandlers(service)

	ticket := &domain.Ticket{TicketId: 101, VehicleNumber: "UP16AB1234", SlotId: 1, EntryTime: time.Now().Add(-1 * time.Hour)}
	slot := &domain.Slot{SlotId: 1, SlotType: "car", IsFree: false}

	ticketRepo.On("FindTicketByVehicleNumber", ticket.VehicleNumber).Return(ticket, nil)
	slotRepo.On("FindSlotByID", ticket.SlotId).Return(slot, nil)
	slotRepo.On("FindSlotTypebyID", ticket.SlotId).Return("car", nil)
	slotRepo.On("UpdateSlot", mock.Anything).Return(nil)
	ticketRepo.On("DeleteTicket", ticket.TicketId).Return(nil)

	body := []byte(`{"vehiclenumber":"UP16AB1234"}`)
	req := httptest.NewRequest(http.MethodPost, "/unpark", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.UnparkVehicleRequest(w, req)
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestAddSlot(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)
	handler := NewHandlers(service)

	slot := domain.Slot{SlotId: 1, SlotType: "car", IsFree: true}
	slotRepo.On("SaveSlot", slot).Return(nil)

	body, _ := json.Marshal(slot)
	req := httptest.NewRequest(http.MethodPost, "/add-slot", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.AddSlot(w, req)
	res := w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestGetAvailableSlots(t *testing.T) {
	slotRepo := new(MockSlotRepo)
	ticketRepo := new(MockTicketRepo)
	service := parking.NewParkingService(slotRepo, ticketRepo)
	handler := NewHandlers(service)

	slots := []domain.Slot{{SlotId: 1, SlotType: "car", IsFree: true}}
	slotRepo.On("ListAvailableSlots").Return(slots, nil)

	req := httptest.NewRequest(http.MethodGet, "/available-slots", nil)
	w := httptest.NewRecorder()

	handler.GetAvailableSlots(w, req)
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
