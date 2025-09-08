package mysql

import (
	"database/sql"
	"errors"
	"parkingSlotManagement/internals/core/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveSlot(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSlotRepo(db)

	tests := []struct {
		name          string
		slot          domain.Slot
		mockBehavior  func()
		expectedError bool
	}{
		{
			name: "successfully saves slot",
			slot: domain.Slot{SlotId: 1, SlotType: "car", IsFree: true},
			mockBehavior: func() {
				mock.ExpectExec("INSERT INTO slots").
					WithArgs(1, "car", true).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "fails to save slot due to DB error",
			slot: domain.Slot{SlotId: 2, SlotType: "bike", IsFree: false},
			mockBehavior: func() {
				mock.ExpectExec("INSERT INTO slots").
					WithArgs(2, "bike", false).
					WillReturnError(errors.New("error inserting slot"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			err := repo.SaveSlot(tt.slot)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateSlot(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	SlotRepo := NewSlotRepo(db)
	tests := []struct {
		name          string
		slot          domain.Slot
		mockFunc      func()
		expectedError bool
	}{
		{
			name: "successfully update slot",
			slot: domain.Slot{SlotId: 1, SlotType: "car", IsFree: false},
			mockFunc: func() {
				mock.ExpectExec(`(?i)UPDATE\s+slots\s+SET\s+slottype=\?,\s*isfree=\?\s+WHERE\s+slotid=\?`).
					WithArgs("car", false, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))

			},
			expectedError: false,
		},
		{
			name: "fail to update slot in DB",
			slot: domain.Slot{SlotId: 1, SlotType: "car", IsFree: false},
			mockFunc: func() {
				mock.ExpectExec(`(?i)UPDATE\s+slots\s+SET\s+slottype=\?,\s*isfree=\?\s+WHERE\s+slotid=\?`).
					WithArgs("car", false, 1).
					WillReturnError(errors.New("error updating slot"))

			},
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := SlotRepo.UpdateSlot(&tt.slot)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListAvailableSlots(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	SlotRepo := NewSlotRepo(db)
	tests := []struct {
		name          string
		mockFunc      func()
		expectedSlots []domain.Slot
		expectedError bool
	}{
		{
			name: "successfully get available slots",
			mockFunc: func() {
				mock.ExpectQuery("SELECT slotid,slottype,isfree FROM slots WHERE isfree=true").
					WillReturnRows(sqlmock.NewRows([]string{"slotid", "slottype", "isfree"}).
						AddRow(1, "car", true).
						AddRow(2, "bike", true))
			},
			expectedSlots: []domain.Slot{
				{SlotId: 1, SlotType: "car", IsFree: true},
				{SlotId: 2, SlotType: "bike", IsFree: true},
			},
			expectedError: false,
		},
		{
			name: "failed to  get available slots",
			mockFunc: func() {
				mock.ExpectQuery("SELECT slotid,slottype,isfree FROM slots WHERE isfree=true").
					WillReturnError(errors.New("error fetching slots"))
			},
			expectedSlots: nil,
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			_, err := SlotRepo.ListAvailableSlots()
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}
}

func TestFindSlotByType(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	SlotRepo := NewSlotRepo(db)
	tests := []struct {
		name          string
		slotType      string
		mockFunc      func(slotType string)
		expectedSlots []domain.Slot
		expectedError bool
	}{
		{
			name:     "successfully get slots by type",
			slotType: "car",
			mockFunc: func(slotType string) {
				mock.ExpectQuery(`(?i)SELECT\s+slotid,\s*slottype,\s*isfree\s+FROM\s+slots\s+WHERE\s+slottype=\?\s+AND\s+isfree=true`).
					WithArgs(slotType).
					WillReturnRows(sqlmock.NewRows([]string{"slotid", "slottype", "isfree"}).
						AddRow(1, "car", true).
						AddRow(2, "bike", true))

			},
			expectedSlots: []domain.Slot{
				{SlotId: 1, SlotType: "car", IsFree: true},
			},
			expectedError: false,
		},
		{
			name:     "failed get slots by type",
			slotType: "car",
			mockFunc: func(slotType string) {
				mock.ExpectQuery(`(?i)SELECT\s+slotid,\s*slottype,\s*isfree\s+FROM\s+slots\s+WHERE\s+slottype=\?\s+AND\s+isfree=true`).
					WithArgs(slotType).
					WillReturnError(errors.New("error fetching slot by type"))
			},
			expectedSlots: nil,
			expectedError: true,
		},
	}
	// abcd
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.slotType)

			_, err := SlotRepo.FindSlotByType(tt.slotType)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}
}
func TestFindSlotTypeByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSlotRepo(db)

	tests := []struct {
		name          string
		slotID        int
		mockFunc      func()
		expectedType  string
		expectedError bool
	}{
		{
			name:   "successfully fetch slot type",
			slotID: 1,
			mockFunc: func() {
				mock.ExpectQuery(`(?i)SELECT\s+slottype\s+from\s+slots\s+WHERE\s+slotid=\?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"slottype"}).AddRow("car"))
			},
			expectedType:  "car",
			expectedError: false,
		},
		{
			name:   "slot not found",
			slotID: 2,
			mockFunc: func() {
				mock.ExpectQuery(`(?i)SELECT\s+slottype\s+from\s+slots\s+WHERE\s+slotid=\?`).
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			expectedType:  "",
			expectedError: true,
		},
		{
			name:   "db error",
			slotID: 3,
			mockFunc: func() {
				mock.ExpectQuery(`(?i)SELECT\s+slottype\s+from\s+slots\s+WHERE\s+slotid=\?`).
					WithArgs(3).
					WillReturnError(errors.New("db error"))
			},
			expectedType:  "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			slotType, err := repo.FindSlotTypebyID(tt.slotID)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedType, slotType)
			}
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}
}

func TestFindSlotByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSlotRepo(db)

	tests := []struct {
		name          string
		slotID        int
		mockFunc      func()
		expectedSlot  *domain.Slot
		expectedError bool
	}{
		{
			name:   "successfully fetch slot by ID",
			slotID: 1,
			mockFunc: func() {
				mock.ExpectQuery(`(?i)SELECT\s+slotid,\s*slottype,\s*isfree\s+FROM\s+slots\s+WHERE\s+slotid\s*=\s*\?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"slotid", "slottype", "isfree"}).
						AddRow(1, "car", true))
			},
			expectedSlot:  &domain.Slot{SlotId: 1, SlotType: "car", IsFree: true},
			expectedError: false,
		},
		{
			name:   "slot not found",
			slotID: 2,
			mockFunc: func() {
				mock.ExpectQuery(`(?i)SELECT\s+slotid,\s*slottype,\s*isfree\s+FROM\s+slots\s+WHERE\s+slotid\s*=\s*\?`).
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			expectedSlot:  nil,
			expectedError: true,
		},
		{
			name:   "db error",
			slotID: 3,
			mockFunc: func() {
				mock.ExpectQuery(`(?i)SELECT\s+slotid,\s*slottype,\s*isfree\s+FROM\s+slots\s+WHERE\s+slotid\s*=\s*\?`).
					WithArgs(3).
					WillReturnError(errors.New("db error"))
			},
			expectedSlot:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			slot, err := repo.FindSlotByID(tt.slotID)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSlot, slot)
			}
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}
}
