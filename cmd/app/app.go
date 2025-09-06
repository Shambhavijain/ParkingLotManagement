package app

import (
	"log"
	"net/http"
	"parkingSlotManagement/internals/adapters/repositories/mysql"
	"parkingSlotManagement/internals/adapters/requestHandlers"
	"parkingSlotManagement/internals/core/services/parking"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func Start() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error Loading .env file")
	}
	database := mysql.GetInstance()

	SlotRepo := mysql.NewSlotRepo(database)
	TicketRepo := mysql.NewTicketRepo(database)

	//InMemmory
	// SlotRepo := inmemmory.NewSlotInMemmory()
	// TicketRepo := inmemmory.NewTicketInMemmory()

	ParkingService := parking.NewParkingService(SlotRepo, TicketRepo)
	handler := requestHandlers.NewHandlers(ParkingService)

	r := mux.NewRouter()

	r.HandleFunc("/ParkVehicle", handler.ParkVehicleRequest).Methods(http.MethodPost)
	r.HandleFunc("/UnparkVehicle", handler.UnparkVehicleRequest).Methods(http.MethodPost)
	r.HandleFunc("/AddSlot", handler.AddSlot).Methods(http.MethodPost)
	log.Println("Server running on:8080")
	http.ListenAndServe(":8080", r)
}
