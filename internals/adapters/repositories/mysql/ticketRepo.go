package mysql

import (
	"database/sql"
	"parkingSlotManagement/internals/core/domain"
	"time"
)

type TicketRepo struct {
	db *sql.DB
}

func NewTicketRepo(db *sql.DB) *TicketRepo {
	return &TicketRepo{db: db}
}
func (t *TicketRepo) SaveTicket(ticket domain.Ticket) error {
	_, err := t.db.Exec("INSERT INTO  tickets (ticketid,vehiclenumber,entrytime,slotid)VALUE (?,?,?,?)",
		ticket.TicketId, ticket.VehicleNumber, ticket.EntryTime, ticket.SlotId)

	return err

}
func (t *TicketRepo) DeleteTicket(ticketid int64) error {
	_, err := t.db.Exec("DELETE FROM tickets WHERE ticketid=?", ticketid)

	if err != nil {
		return err
	}
	return nil
}

func (t *TicketRepo) FindTicketByVehicleNumber(Vehiclenumber string) (*domain.Ticket, error) {
	var Ticket domain.Ticket
	var entryTimeStr string

	row := t.db.QueryRow("SELECT ticketid, vehiclenumber, entrytime, slotid FROM tickets WHERE vehiclenumber = ?", Vehiclenumber)
	err := row.Scan(&Ticket.TicketId, &Ticket.VehicleNumber, &entryTimeStr, &Ticket.SlotId)
	if err != nil {
		return nil, err
	}

	Ticket.EntryTime, err = time.Parse("2006-01-02 15:04:05", entryTimeStr)
	if err != nil {
		return nil, err
	}

	return &Ticket, nil

}
