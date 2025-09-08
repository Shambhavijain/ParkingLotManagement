package mysql

import (
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
				mock.ExpectExec("UPDATE slots SET slottype=?, isfree=? WHERE slotid=?").
					WithArgs("car", false, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: false,
		},
		{
			name: "fail to update slot in DB",
			slot: domain.Slot{SlotId: 1, SlotType: "car", IsFree: false},
			mockFunc: func() {
				mock.ExpectExec("UPDATE slots SET slottype=?, isfree=? WHERE slotid=?").
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
				mock.ExpectQuery("SELECT slotid,slottype,isfree FROM slots where slottype=? AND isfree=true").
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
				mock.ExpectQuery("SELECT slotid,slottype,isfree FROM slots where slottype=? AND isfree=true").
					WithArgs(slotType).
					WillReturnError(errors.New("error fetching slot by type"))
			},
			expectedSlots: nil,
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.slotType)
		})
		_, err := SlotRepo.FindSlotByType(tt.slotType)
		if tt.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Nil(t, mock.ExpectationsWereMet())
	}
}
