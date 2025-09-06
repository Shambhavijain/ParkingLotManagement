package ports

import "parkingSlotManagement/internals/core/domain"

type TicketRepository interface {
	SaveTicket(ticket domain.Ticket) error
	FindTicketByVehicleNumber(vehiclenumber string) (*domain.Ticket, error)
	DeleteTicket(ticketid int64) error
}
