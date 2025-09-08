package mysql

import (
	"errors"
	"fmt"
)

var (
	ErrSlotNotFound     = errors.New("slot not found")
	ErrTicketNotFound   = errors.New("ticket not found")
	ErrSlotNotFoundByID = errors.New(" not found slot by  slot ID")
	ErrInvalidSlotType  = errors.New("invalid slot type")
	ErrDBQueryFailed    = errors.New("database query failed")
)

func Wrap(content string, err error) error {
	if err != nil {

		return fmt.Errorf("%s: %w", content, err)
	}
	return nil
}
