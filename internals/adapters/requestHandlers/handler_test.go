package requestHandlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"parkingSlotManagement/internals/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockParkingService struct {
	mock.Mock
}

func (m *MockParkingService) ParkVehicle(vehicle domain.Vehicle) (*domain.Ticket, error) {
	args := m.Called(vehicle)
	return args.Get(0).(*domain.Ticket), args.Error(1)
}

func (m *MockParkingService) UnparkVehicle(vehicleNumber string) (float64, error) {
	args := m.Called(vehicleNumber)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockParkingService) AddSlot(slot domain.Slot) error {
	args := m.Called(slot)
	return args.Error(0)
}

func (m *MockParkingService) GetAvailableSlots() ([]domain.Slot, error) {
	args := m.Called()
	return args.Get(0).([]domain.Slot), args.Error(1)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenStr string) (*domain.Admin, error) {
	args := m.Called(tokenStr)
	return args.Get(0).(*domain.Admin), args.Error(1)
}

func TestParkVehicleRequest(t *testing.T) {
	mockService := new(MockParkingService)
	authService := new(MockAuthService)
	handler := NewHandlers(mockService, authService)

	vehicle := domain.Vehicle{
		VehicleNumber: "UP16AB1234",
		VehicleType:   "car",
	}

	ticket := &domain.Ticket{
		TicketId:      123456,
		VehicleNumber: "UP16AB1234",
		SlotId:        1,
	}

	mockService.On("ParkVehicle", vehicle).Return(ticket, nil)

	body, _ := json.Marshal(vehicle)
	req := httptest.NewRequest(http.MethodPost, "/ParkVehicleRequest", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ParkVehicleRequest(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}
func TestParkVehicleRequest_AlreadyParked(t *testing.T) {
	mockService := new(MockParkingService)
	authService := new(MockAuthService)
	handler := NewHandlers(mockService, authService)

	vehicle := domain.Vehicle{
		VehicleNumber: "UP16AB1234",
		VehicleType:   "car",
	}

	mockService.On("ParkVehicle", vehicle).Return((*domain.Ticket)(nil), errors.New("vehicle already parked"))

	body, _ := json.Marshal(vehicle)
	req := httptest.NewRequest(http.MethodPost, "/ParkVehicleRequest", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ParkVehicleRequest(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "vehicle already parked")
}
func TestParkVehicleRequest_InvalidJSON(t *testing.T) {
	mockService := new(MockParkingService)
	authService := new(MockAuthService)
	handler := NewHandlers(mockService, authService)

	req := httptest.NewRequest(http.MethodPost, "/ParkVehicleRequest", bytes.NewReader([]byte("invalid-json")))
	w := httptest.NewRecorder()

	handler.ParkVehicleRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Body Request")
}

func TestUnparkVehicleRequest(t *testing.T) {
	mockService := new(MockParkingService)
	authService := new(MockAuthService)
	handler := NewHandlers(mockService, authService)

	vehicleNumber := "UP16AB1234"
	expectedFee := 120.0

	mockService.On("UnparkVehicle", vehicleNumber).Return(expectedFee, nil)

	body, _ := json.Marshal(map[string]string{
		"vehiclenumber": vehicleNumber,
	})

	req := httptest.NewRequest(http.MethodPost, "/UnparkVehicleRequest", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.UnparkVehicleRequest(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, vehicleNumber, response["vehiclenumber"])
	assert.Equal(t, 120.0, response["fee"])
	assert.Equal(t, "Successfully Unpark The Vehicle", response["message"])

	mockService.AssertExpectations(t)
}
func TestUnparkVehicleRequest_InvalidJSON(t *testing.T) {
	mockService := new(MockParkingService)
	authService := new(MockAuthService)
	handler := NewHandlers(mockService, authService)

	req := httptest.NewRequest(http.MethodPost, "/UnparkVehicleRequest", bytes.NewReader([]byte("invalid-json")))
	w := httptest.NewRecorder()

	handler.UnparkVehicleRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Body Request")
}

func TestAddSlot(t *testing.T) {
	mockService := new(MockParkingService)
	authService := new(MockAuthService)
	handler := NewHandlers(mockService, authService)

	slot := domain.Slot{
		SlotId:   1,
		SlotType: "car",
		IsFree:   true,
	}

	mockService.On("AddSlot", slot).Return(nil)

	body, _ := json.Marshal(slot)
	req := httptest.NewRequest(http.MethodPost, "/AddSlot", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.AddSlot(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "Slot added successfully", w.Body.String())

	mockService.AssertExpectations(t)
}
func TestUnparkVehicleRequest_TicketNotFound(t *testing.T) {
	mockService := new(MockParkingService)
	authService := new(MockAuthService)

	handler := NewHandlers(mockService, authService)

	vehicleNumber := "UP16AB1234"

	mockService.On("UnparkVehicle", vehicleNumber).Return(0.0, errors.New("ticket not found"))

	body, _ := json.Marshal(map[string]string{
		"vehiclenumber": vehicleNumber,
	})
	req := httptest.NewRequest(http.MethodPost, "/unpark", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.UnparkVehicleRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "ticket not found")
}

func TestGetAvailableSlots(t *testing.T) {
	mockService := new(MockParkingService)
	authService := new(MockAuthService)

	handler := NewHandlers(mockService, authService)

	slots := []domain.Slot{
		{SlotId: 1, SlotType: "car", IsFree: true},
		{SlotId: 2, SlotType: "bike", IsFree: true},
	}

	mockService.On("GetAvailableSlots").Return(slots, nil)

	req := httptest.NewRequest(http.MethodGet, "/available-slots", nil)
	w := httptest.NewRecorder()

	handler.GetAvailableSlots(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []domain.Slot
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, slots, response)

	mockService.AssertExpectations(t)
}

func TestLoginHandler(t *testing.T) {
	mockService := new(MockParkingService)
	authService := new(MockAuthService)

	handler := NewHandlers(mockService, authService)

	creds := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	expectedToken := "mocked.jwt.token"

	authService.On("Login", creds["username"], creds["password"]).Return(expectedToken, nil)

	body, _ := json.Marshal(creds)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.LoginHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, expectedToken, response["token"])
	assert.Equal(t, "Login successful", response["message"])

	authService.AssertExpectations(t)
}
