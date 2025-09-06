package domain

type Slot struct {
	SlotId   int    `json:"slotid"`
	SlotType string `json:"slottype"`
	IsFree   bool   `json:"isfree"`
}
