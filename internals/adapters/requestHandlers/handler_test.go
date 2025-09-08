package requestHandlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"parkingSlotManagement/internals/adapters/repositories/inmemmory"
	"parkingSlotManagement/internals/core/domain"
	"parkingSlotManagement/internals/core/services/parking"
	"strings"
	"testing"
)

func TestAddSlot(t *testing.T) {
	slotRepo := inmemmory.NewSlotInMemmory()
	ticketRepo := inmemmory.NewTicketInMemmory()

	service := parking.NewParkingService(slotRepo, ticketRepo)
	h := NewHandlers(service)

	Slot := domain.Slot{
		SlotId:   1,
		SlotType: "car",
		IsFree:   true,
	}

	body, err := json.Marshal(Slot)
	if err != nil {
		t.Fatalf("Failed to marshal slot: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/AddSlot", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	h.AddSlot(resp, req)
	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201 Created , got %d", resp.Code)
	}

	expected1 := "Slot added successfully"
	if strings.TrimSpace(resp.Body.String()) != expected1 {
		t.Errorf("Valid Slot: Expected body %q, got %q", expected1, resp.Body.String())
	}

	req2 := httptest.NewRequest(http.MethodPost, "/AddSlot", strings.NewReader("{invalid json"))
	req2.Header.Set("Content-Type", "application/json")
	resp2 := httptest.NewRecorder()

	h.AddSlot(resp2, req2)

	if resp2.Code != http.StatusBadRequest {
		t.Errorf("Invalid JSON: Expected status 400 Bad Request, got %d", resp2.Code)
	}

	expected2 := "Invalid Body Request"
	if !strings.Contains(resp2.Body.String(), expected2) {
		t.Errorf("Invalid JSON: Expected error message %q, got %q", expected2, resp2.Body.String())
	}

}

func TestGetAvailableSlot(t *testing.T) {
	slotRepo := inmemmory.NewSlotInMemmory()
	ticketRepo := inmemmory.NewTicketInMemmory()

	service := parking.NewParkingService(slotRepo, ticketRepo)
	h := NewHandlers(service)

	Slot := domain.Slot{
		SlotId:   1,
		SlotType: "car",
		IsFree:   true,
	}

	body, err := json.Marshal(Slot)
	if err != nil {
		t.Fatalf("Failed to marshal slot: %v", err)
	}
	slotRepo.SaveSlot(Slot)
	req := httptest.NewRequest(http.MethodPost, "/GetAvailableSlot", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	h.GetAvailableSlots(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200 Created , got %d", resp.Code)
	}

	var slots []domain.Slot
	if err := json.NewDecoder(resp.Body).Decode(&slots); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

}
func TestParkVehicleRequest(t *testing.T) {

	slotRepo := inmemmory.NewSlotInMemmory()
	ticketRepo := inmemmory.NewTicketInMemmory()
	service := parking.NewParkingService(slotRepo, ticketRepo)
	h := NewHandlers(service)

	slot := domain.Slot{
		SlotId:   1,
		SlotType: "car",
		IsFree:   true,
	}
	slotRepo.SaveSlot(slot)

	vehicle := domain.Vehicle{
		VehicleNumber: "UP16AB1234",
		VehicleType:   "car",
	}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(vehicle)

	req := httptest.NewRequest(http.MethodPost, "/ParkVehicle", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	h.ParkVehicleRequest(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", resp.Code)
	}

	var ticket domain.Ticket
	if err := json.NewDecoder(resp.Body).Decode(&ticket); err != nil {
		t.Fatalf("Failed to decode ticket: %v", err)
	}

	req2 := httptest.NewRequest(http.MethodPost, "/ParkVehicle", strings.NewReader("{invalid json"))
	req2.Header.Set("Content-Type", "application/json")
	resp2 := httptest.NewRecorder()

	h.ParkVehicleRequest(resp2, req2)

	if resp2.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 Internal Server Error for invalid JSON, got %d", resp2.Code)
	}

	if !strings.Contains(resp2.Body.String(), "Invalid Body Request") {
		t.Errorf("Expected error message 'Invalid Body Request', got %q", resp2.Body.String())
	}

	failVehicle := domain.Vehicle{
		VehicleNumber: "FAIL123",
		VehicleType:   "car",
	}
	body3 := new(bytes.Buffer)
	json.NewEncoder(body3).Encode(failVehicle)

	req3 := httptest.NewRequest(http.MethodPost, "/ParkVehicle", body3)
	req3.Header.Set("Content-Type", "application/json")
	resp3 := httptest.NewRecorder()

	h.ParkVehicleRequest(resp3, req3)

	if resp3.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 Internal Server Error for service error, got %d", resp3.Code)
	}

	if !strings.Contains(resp3.Body.String(), "no available slots for this vehicle type\n") {
		t.Errorf("Expected error message 'parking failed', got %q", resp3.Body.String())
	}
}

func TestUnparkVehicleRequest(t *testing.T) {

	slotRepo := inmemmory.NewSlotInMemmory()
	ticketRepo := inmemmory.NewTicketInMemmory()
	service := parking.NewParkingService(slotRepo, ticketRepo)
	h := NewHandlers(service)

	slot := domain.Slot{
		SlotId:   1,
		SlotType: "car",
		IsFree:   true,
	}
	slotRepo.SaveSlot(slot)

	vehicle := domain.Vehicle{
		VehicleNumber: "UP16AB1234",
		VehicleType:   "car",
	}
	_, err := service.ParkVehicle(vehicle)
	if err != nil {
		t.Fatalf("Failed to park vehicle: %v", err)
	}

	unparkReq := map[string]string{
		"vehiclenumber": "UP16AB1234",
	}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(unparkReq)

	req := httptest.NewRequest(http.MethodPost, "/UnparkVehicle", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	h.UnparkVehicleRequest(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", resp.Code)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["vehiclenumber"] != "UP16AB1234" {
		t.Errorf("Expected vehiclenumber 'UP16AB1234', got %v", result["vehiclenumber"])
	}

	if result["message"] != "Successfully Unpark The Vehicle" {
		t.Errorf("Expected success message, got %v", result["message"])
	}

	req2 := httptest.NewRequest(http.MethodPost, "/UnparkVehicle", strings.NewReader("{invalid json"))
	req2.Header.Set("Content-Type", "application/json")
	resp2 := httptest.NewRecorder()

	h.UnparkVehicleRequest(resp2, req2)

	if resp2.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 Internal Server Error for invalid JSON, got %d", resp2.Code)
	}

	if !strings.Contains(resp2.Body.String(), "Invalid Body Request") {
		t.Errorf("Expected error message 'Invalid Body Request', got %q", resp2.Body.String())
	}

	unparkReq3 := map[string]string{
		"vehiclenumber": "NOTEXIST123",
	}
	body3 := new(bytes.Buffer)
	json.NewEncoder(body3).Encode(unparkReq3)

	req3 := httptest.NewRequest(http.MethodPost, "/UnparkVehicle", body3)
	req3.Header.Set("Content-Type", "application/json")
	resp3 := httptest.NewRecorder()

	h.UnparkVehicleRequest(resp3, req3)

	if resp3.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 Internal Server Error for unknown vehicle, got %d", resp3.Code)
	}

	if !strings.Contains(resp3.Body.String(), "ticket of this vehiclenumber not found") {
		t.Errorf("Expected error message 'ticket not found', got %q", resp3.Body.String())
	}
}
