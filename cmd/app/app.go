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
	
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	
	db := mysql.GetInstance()

	
	slotRepo := mysql.NewSlotRepo(db)
	ticketRepo := mysql.NewTicketRepo(db)
	userRepo := mysql.NewMySQLUserRepository(db) 

	
	parkingService := parking.NewParkingService(slotRepo, ticketRepo)
	authService := auth.NewAuthService(userRepo)

	
	handler := requestHandlers.NewHandlers(parkingService)

	
	r := mux.NewRouter()

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handler.LoginHandler(w, r, authService)
	}).Methods(http.MethodPost)

	r.Handle("/ParkVehicle", middleware.AuthMiddleware(authService)(http.HandlerFunc(handler.ParkVehicleRequest))).Methods(http.MethodPost)
	r.Handle("/UnparkVehicle", middleware.AuthMiddleware(authService)(http.HandlerFunc(handler.UnparkVehicleRequest))).Methods(http.MethodPost)
	r.Handle("/GetAvailableSlots", middleware.AuthMiddleware(authService)(http.HandlerFunc(handler.GetAvailableSlots))).Methods(http.MethodPost)
	r.Handle("/AddSlot", middleware.AuthMiddleware(authService)(http.HandlerFunc(handler.AddSlot))).Methods(http.MethodPost)

	
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
