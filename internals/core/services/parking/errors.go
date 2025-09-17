package parking

import (
	"errors"
	"fmt"
)
//type AppErr error
var (
	ErrTicketNotFound       = errors.New("ticket of this vehicle number not found")
	ErrSlotNotFound         = errors.New("failed to fetch slot")
	ErrSlotUpdateFailed     = errors.New("failed to update slot status")
	ErrTicketDeleteFailed   = errors.New("ticket can't be deleted")
	ErrFeeCalculationFailed = errors.New("unable to calculate fee")
	ErrInvalidVehicleType   = errors.New("invalid vehicle type")
	ErrSlotSaveFailed       = errors.New("error inserting slot")
	ErrSlotListFailed       = errors.New("error fetching available slots")
	ErrTicketSaveFailed     = errors.New("failed to save ticket to database")
	ErrExistingTicketCheck  = errors.New("error checking existing ticket")
	ErrSlotFetchByType      = errors.New("failed to fetch slots by type ")
	ErrVehicleAlreadyParked = errors.New("vehicle has been already parked")
)

func Wrap(content string, err error) error {
	if err != nil {

		return fmt.Errorf("%s: %w", content, err)
	}
	return nil
}

// func (e AppErr) ToHTTPError() 

// func AppErrToHttpErr(err error) httpError {

// }
