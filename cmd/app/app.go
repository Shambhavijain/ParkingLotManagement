package app

import (
	"log"
	"net/http"
	"parkingSlotManagement/internals/adapters/repositories/mysql"
	"parkingSlotManagement/internals/adapters/requestHandlers"
	"parkingSlotManagement/internals/adapters/requestHandlers/middleware"
	"parkingSlotManagement/internals/core/services/auth"
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
	AuthService := auth.NewAuthService()
	handler := requestHandlers.NewHandlers(ParkingService)

	loginHandler := requestHandlers.LoginHandler(AuthService)

	r := mux.NewRouter()

	r.HandleFunc("/login", loginHandler).Methods(http.MethodPost)


	r.HandleFunc("/ParkVehicle", middleware.AuthMiddleware(handler.ParkVehicleRequest, AuthService)).Methods(http.MethodPost)
	r.HandleFunc("/UnparkVehicle", middleware.AuthMiddleware(handler.UnparkVehicleRequest, AuthService)).Methods(http.MethodPost)
	r.HandleFunc("/AddSlot", middleware.AuthMiddleware(handler.AddSlot, AuthService)).Methods(http.MethodPost)
	r.HandleFunc("/GetAvailableSlots", middleware.AuthMiddleware(handler.GetAvailableSlots, AuthService)).Methods(http.MethodPost)

	log.Println("Server running on:8080")
	http.ListenAndServe(":8080", r)
}
