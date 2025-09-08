package main

import (
	"bufio"
	"fmt"
	"os"
	"parkingSlotManagement/internals/adapters/repositories/inmemmory"
	"parkingSlotManagement/internals/core/domain"
	"parkingSlotManagement/internals/core/services/parking"
	"strconv"
	"strings"
)

func main() {

	slotRepo := inmemmory.NewSlotInMemmory()
	ticketRepo := inmemmory.NewTicketInMemmory()

	service := parking.NewParkingService(slotRepo, ticketRepo)

	// Pre-populate slots
	slotRepo.SaveSlot(domain.Slot{SlotId: 1, SlotType: "car", IsFree: true})
	slotRepo.SaveSlot(domain.Slot{SlotId: 2, SlotType: "car", IsFree: true})
	slotRepo.SaveSlot(domain.Slot{SlotId: 3, SlotType: "bike", IsFree: true})
	slotRepo.SaveSlot(domain.Slot{SlotId: 4, SlotType: "bike", IsFree: true})

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Parking Lot Management System 🅿️")

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
				fmt.Printf("Vehicle parked successfully. Ticket ID: %d\n", ticket.TicketId)
			}

		case "2":
			fmt.Print("Enter vehicle number: ")
			number, _ := reader.ReadString('\n')
			number = strings.TrimSpace(number)

			fee, err := service.UnparkVehicle(number)
			if err != nil {
				fmt.Printf(" Error: %v\n", err)
			} else {
				fmt.Printf(" Vehicle unparked. Fee: ₹%.2f\n", fee)
			}

		case "3":
			slots, err := service.GetAvailableSlots()
			if err != nil {
				fmt.Printf("Error fetching slots: %v\n", err)
			} else {
				fmt.Println(" Available Slots:")
				for _, slot := range slots {
					fmt.Printf("Slot ID: %d | Type: %s\n", slot.SlotId, slot.SlotType)
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
