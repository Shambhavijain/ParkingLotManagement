package mysql

import (
	"errors"
	"parkingSlotManagement/internals/core/domain"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveTicket(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTicketRepo(db)

	//entryTime := time.Date(2025, 9, 8, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		ticket        domain.Ticket
		mockFunc      func(ticket domain.Ticket)
		expectedError bool
	}{
		{
			name: "successfully save ticket",
			ticket: domain.Ticket{
				TicketId:      1,
				VehicleNumber: "UP16AB1234",
				EntryTime:     time.Date(2025, 9, 8, 10, 0, 0, 0, time.UTC),
				SlotId:        1,
			},
			mockFunc: func(ticket domain.Ticket) {

				mock.ExpectExec(`(?i)INSERT\s+INTO\s+tickets\s*\(ticketid,vehiclenumber,entrytime,slotid\)\s*VALUES\s*\(\?,\?,\?,\?\)`).
					WithArgs(ticket.TicketId, ticket.VehicleNumber, sqlmock.AnyArg(), ticket.SlotId).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			expectedError: false,
		},
		{
			name: "fail to save ticket",
			ticket: domain.Ticket{
				TicketId:      2,
				VehicleNumber: "UP16XY5678",
				EntryTime:     time.Date(2025, 9, 8, 10, 0, 0, 0, time.UTC),
				SlotId:        2,
			},
			mockFunc: func(ticket domain.Ticket) {
				mock.ExpectExec(`(?i)INSERT\s+INTO\s+tickets\s*\(ticketid,vehiclenumber,entrytime,slotid\)\s*VALUES\s*\(\?,\?,\?,\?\)`).
					WithArgs(ticket.TicketId, ticket.VehicleNumber, sqlmock.AnyArg(), ticket.SlotId).
					WillReturnError(errors.New("insert failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.ticket)
			err := repo.SaveTicket(tt.ticket)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}

}

func TestDeleteTicket(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTicketRepo(db)

	tests := []struct {
		name          string
		ticketID      int64
		mockFunc      func()
		expectedError bool
	}{
		{
			name:     "successfully delete ticket",
			ticketID: 1,
			mockFunc: func() {
				mock.ExpectExec(`(?i)DELETE\s+FROM\s+tickets\s+WHERE\s+ticketid=\?`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: false,
		},
		{
			name:     "fail to delete ticket",
			ticketID: 2,
			mockFunc: func() {
				mock.ExpectExec(`(?i)DELETE\s+FROM\s+tickets\s+WHERE\s+ticketid=\?`).
					WithArgs(2).
					WillReturnError(errors.New("delete error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := repo.DeleteTicket(tt.ticketID)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}
}
func TestFindTicketByVehicleNumber(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTicketRepo(db)

	tests := []struct {
		name           string
		vehicleNumber  string
		mockFunc       func()
		expectedTicket *domain.Ticket
		expectedError  bool
	}{
		{
			name:          "successfully find ticket",
			vehicleNumber: "UP16AB1234",
			mockFunc: func() {
				mock.ExpectQuery(`(?i)SELECT\s+ticketid,\s*vehiclenumber,\s*entrytime,\s*slotid\s+FROM\s+tickets\s+WHERE\s+vehiclenumber\s*=\s*\?`).
					WithArgs("UP16AB1234").
					WillReturnRows(sqlmock.NewRows([]string{"ticketid", "vehiclenumber", "entrytime", "slotid"}).
						AddRow(1, "UP16AB1234", "2025-09-08 10:00:00", 101))
			},
			expectedTicket: &domain.Ticket{
				TicketId:      1,
				VehicleNumber: "UP16AB1234",
				EntryTime:     time.Date(2025, 9, 8, 10, 0, 0, 0, time.UTC),
				SlotId:        101,
			},
			expectedError: false,
		},
		{
			name:          "fail to find ticket",
			vehicleNumber: "UP16XY5678",
			mockFunc: func() {
				mock.ExpectQuery(`(?i)SELECT\s+ticketid,\s*vehiclenumber,\s*entrytime,\s*slotid\s+FROM\s+tickets\s+WHERE\s+vehiclenumber\s*=\s*\?`).
					WithArgs("UP16XY5678").
					WillReturnError(errors.New("query error"))
			},
			expectedTicket: nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			ticket, err := repo.FindTicketByVehicleNumber(tt.vehicleNumber)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTicket, ticket)
			}
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}
}
