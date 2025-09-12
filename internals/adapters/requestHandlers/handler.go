package requestHandlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"parkingSlotManagement/internals/core/domain"
	"parkingSlotManagement/internals/core/services/auth"
	"parkingSlotManagement/internals/core/services/parking"
)

type Handlers struct {
	service *parking.ParkingService
}

func NewHandlers(service *parking.ParkingService) *Handlers {
	return &Handlers{
		service: service,
	}
}
func (h *Handlers) ParkVehicleRequest(w http.ResponseWriter, r *http.Request) {
	var vehicle domain.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		http.Error(w, "Invalid Body Request", http.StatusInternalServerError)
		return
	}
	fmt.Println(vehicle)
	ticket, err := h.service.ParkVehicle(vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ticket)

}
func (h *Handlers) UnparkVehicleRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Vehiclenumber string `json:"vehiclenumber"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Body Request", http.StatusInternalServerError)
		return
	}
	fee, err := h.service.UnparkVehicle(req.Vehiclenumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"vehiclenumber": req.Vehiclenumber,
		"fee":           math.Round(fee*100) / 100,
		"message":       "Successfully Unpark The Vehicle",
	})

}
func (h *Handlers) AddSlot(w http.ResponseWriter, r *http.Request) {
	var Slot domain.Slot
	if err := json.NewDecoder(r.Body).Decode(&Slot); err != nil {
		http.Error(w, "Invalid Body Request", http.StatusBadRequest)
		return
	}
	err := h.service.AddSlot(Slot)
	if err != nil {
		http.Error(w, "Unableto add slot", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Slot added successfully"))

}

func (h *Handlers) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	slots, err := h.service.GetAvailableSlots()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(slots); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request, authService *auth.AuthService) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid Body Request", http.StatusBadRequest)
		return
	}

	token, err := authService.Login(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token":   token,
		"message": "Login successful",
	})
}
