package domain

import "time"

type Ticket struct {
	TicketId      int64     `json:"ticketid"`
	VehicleNumber string    `json:"vehiclenumber"`
	SlotId        int       `json:"slotid"`
	EntryTime     time.Time `json:"entrytime"`
}
