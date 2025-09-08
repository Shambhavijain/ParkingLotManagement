package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"parkingSlotManagement/internals/adapters/repositories/mysql"
	"parkingSlotManagement/internals/core/domain"
	"parkingSlotManagement/internals/core/services/auth"
	"parkingSlotManagement/internals/core/services/parking"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// slotRepo := inmemmory.NewSlotInMemmory()
	// ticketRepo := inmemmory.NewTicketInMemmory()

	database := mysql.GetInstance()
	slotRepo := mysql.NewSlotRepo(database)
	ticketRepo := mysql.NewTicketRepo(database)

	service := parking.NewParkingService(slotRepo, ticketRepo)

	authService := auth.NewAuthService()

	// slotRepo.SaveSlot(domain.Slot{SlotId: 1, SlotType: "car", IsFree: true})
	// slotRepo.SaveSlot(domain.Slot{SlotId: 2, SlotType: "car", IsFree: true})
	// slotRepo.SaveSlot(domain.Slot{SlotId: 3, SlotType: "bike", IsFree: true})
	// slotRepo.SaveSlot(domain.Slot{SlotId: 4, SlotType: "bike", IsFree: true})

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Parking Lot Management System ")
	time.Sleep(500 * time.Millisecond)

	// adminUsername := os.Getenv("ADMIN_USERNAME")
	// adminPassword := os.Getenv("ADMIN_PASSWORD")

	for {
		fmt.Print("Enter username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("Enter password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		token, err := authService.Login(username, password)
		if err != nil {
			fmt.Println("Invalid credentials. Please try again.")
			continue
		}

		fmt.Println("Login successful!")
		fmt.Printf("Your token: %s\n", token)
		break
	}

	for {
		fmt.Println("\n--- Menu ---")
		fmt.Println("1. Park Vehicle")
		fmt.Println("2. Unpark Vehicle")
		fmt.Println("3. View Available Slots")
		fmt.Println("4. Add Slot")
		fmt.Println("5. Exit")
		fmt.Print("Enter your choice: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("Enter vehicle number: ")
			number, _ := reader.ReadString('\n')
			number = strings.TrimSpace(number)
			time.Sleep(500 * time.Millisecond)
			fmt.Print("Enter vehicle type (car/bike): ")
			vtype, _ := reader.ReadString('\n')
			vtype = strings.TrimSpace(strings.ToLower(vtype))

			if vtype != "car" && vtype != "bike" {
				fmt.Println("Invalid vehicle type. Please enter 'car' or 'bike'.")
				continue
			}

			ticket, err := service.ParkVehicle(domain.Vehicle{
				VehicleNumber: number,
				VehicleType:   vtype,
			})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				time.Sleep(500 * time.Millisecond)
				fmt.Println("Vehicle parked successfully. Ticket details:")
				fmt.Printf("Ticket ID: %d\n", ticket.TicketId)
				fmt.Printf("Vehicle Number: %s\n", ticket.VehicleNumber)
				fmt.Printf("Entry Time: %s\n", ticket.EntryTime.Format("2006-01-02 15:04:05"))
				fmt.Printf("Slot ID: %d\n", ticket.SlotId)

			}

		case "2":
			fmt.Print("Enter vehicle number: ")
			number, _ := reader.ReadString('\n')
			number = strings.TrimSpace(number)

			fee, err := service.UnparkVehicle(number)
			if err != nil {
				fmt.Printf(" Error: %v\n", err)
			} else {
				time.Sleep(500 * time.Millisecond)
				fmt.Printf(" Vehicle unparked. Fee: ₹%.2f\n", fee)
			}

		case "3":
			slots, err := service.GetAvailableSlots()
			if err != nil {
				fmt.Printf("Error fetching slots: %v\n", err)
			} else {
				time.Sleep(500 * time.Millisecond)
				fmt.Println(" Available Slots:")
				for _, slot := range slots {
					fmt.Printf("Slot ID: %d | Type: %s   |IsFree: %v\n", slot.SlotId, slot.SlotType, slot.IsFree)
				}
			}

		case "4":
			fmt.Print("Enter new slot ID (number): ")
			idStr, _ := reader.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			slotID, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid slot ID. Please enter a number.")
				continue
			}

			fmt.Print("Enter slot type (car/bike): ")
			slotType, _ := reader.ReadString('\n')
			slotType = strings.TrimSpace(strings.ToLower(slotType))

			if slotType != "car" && slotType != "bike" {
				fmt.Println("Invalid slot type. Please enter 'car' or 'bike'.")
				continue
			}

			err = service.AddSlot(domain.Slot{
				SlotId:   slotID,
				SlotType: slotType,
				IsFree:   true,
			})
			if err != nil {
				fmt.Printf(" Error adding slot: %v\n", err)
			} else {
				fmt.Println(" Slot added successfully.")
			}

		case "5":
			fmt.Println("Thank you for using the Parking Lot System!")
			return

		default:
			fmt.Println("Invalid choice. Please select a valid option.")
		}
	}
}
