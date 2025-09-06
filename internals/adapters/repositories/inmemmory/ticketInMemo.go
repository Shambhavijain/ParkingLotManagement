package inmemmory

import (
	"fmt"
	"parkingSlotManagement/internals/core/domain"
)

type TicketInMemmory struct {
	Tickets map[int]*domain.Ticket
}

func NewTicketInMemmory() *TicketInMemmory {
	return &TicketInMemmory{Tickets: make(map[int]*domain.Ticket)}
}
func (t *TicketInMemmory) SaveTicket(ticket *domain.Ticket) error {
	t.Tickets[ticket.SlotId] = ticket
	return nil
}
func (t *TicketInMemmory) DeleteTicket(ticketid int64) error {
	_, ok := t.Tickets[int(ticketid)]
	if !ok {
		return fmt.Errorf("ticket for this %d id not exists", ticketid)
	}
	delete(t.Tickets, int(ticketid))
	return nil
}
func (t *TicketInMemmory) FindTicketByVehicleNumber(vehiclenumber string) (*domain.Ticket, error) {
	for _, ticket := range t.Tickets {
		if ticket.VehicleNumber == vehiclenumber {
			return ticket, nil
		}
	}
	return nil, fmt.Errorf("unable to find ticket of %s vehiclenumber", vehiclenumber)
}
