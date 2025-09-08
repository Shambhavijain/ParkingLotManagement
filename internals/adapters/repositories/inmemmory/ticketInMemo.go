package inmemmory

import (
	"database/sql"
	"fmt"
	"parkingSlotManagement/internals/core/domain"
)

type TicketInMemmory struct {
	Tickets map[int64]*domain.Ticket
}

func NewTicketInMemmory() *TicketInMemmory {
	return &TicketInMemmory{Tickets: make(map[int64]*domain.Ticket)}
}
func (t *TicketInMemmory) SaveTicket(ticket domain.Ticket) error {
	t.Tickets[ticket.TicketId] = &ticket
	return nil
}
func (t *TicketInMemmory) DeleteTicket(ticketid int64) error {

	_, ok := t.Tickets[ticketid]
	if !ok {
		return fmt.Errorf("ticket for this %d id not exists", ticketid)
	}
	delete(t.Tickets, ticketid)
	return nil

}
func (t *TicketInMemmory) FindTicketByVehicleNumber(vehiclenumber string) (*domain.Ticket, error) {

	for id := range t.Tickets {
		fmt.Println(id)
	}

	for _, ticket := range t.Tickets {
		if ticket.VehicleNumber == vehiclenumber {
			return ticket, nil
		}
	}
	return nil, sql.ErrNoRows

}
