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
	_, err := t.db.Exec("INSERT INTO  tickets (ticketid,vehiclenumber,entrytime,slotid)VALUES (?,?,?,?)",
		ticket.TicketId, ticket.VehicleNumber, ticket.EntryTime, ticket.SlotId)
	if err != nil {
		return ErrDBQueryFailed
	}
	return nil
}
func (t *TicketRepo) DeleteTicket(ticketid int64) error {
	_, err := t.db.Exec("DELETE FROM tickets WHERE ticketid=?", ticketid)

	if err != nil {
		return ErrDBQueryFailed
	}
	return nil
}

func (t *TicketRepo) FindTicketByVehicleNumber(Vehiclenumber string) (*domain.Ticket, error) {
	var Ticket domain.Ticket
	var entryTimeStr string

	row := t.db.QueryRow("SELECT ticketid, vehiclenumber, entrytime, slotid FROM tickets WHERE vehiclenumber = ?", Vehiclenumber)
	err := row.Scan(&Ticket.TicketId, &Ticket.VehicleNumber, &entryTimeStr, &Ticket.SlotId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No ticket found, not an error
		}
		return nil, ErrDBQueryFailed
	}

	Ticket.EntryTime, err = time.Parse("2006-01-02 15:04:05", entryTimeStr)
	if err != nil {
		return nil, Wrap("error parsing entry time", err)
	}

	return &Ticket, nil

}
