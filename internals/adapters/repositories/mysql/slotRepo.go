package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"parkingSlotManagement/internals/core/domain"
)

type SlotRepo struct {
	db *sql.DB
}

func NewSlotRepo(db *sql.DB) *SlotRepo {
	return &SlotRepo{db: db}
}

func (r *SlotRepo) SaveSlot(slot domain.Slot) error {
	_, err := r.db.Exec("INSERT INTO slots (slotid, slottype, isfree) VALUES (?, ?, ?)",
		slot.SlotId, slot.SlotType, slot.IsFree)

	if err != nil {
		log.Printf("Error inserting slot: %v", err)
	}

	return err
}

func (r *SlotRepo) UpdateSlot(slot *domain.Slot) error {
	res, err := r.db.Exec("UPDATE slots SET slottype=?, isfree=? WHERE slotid=?",
		slot.SlotType, slot.IsFree, slot.SlotId)
	if err != nil {
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		fmt.Printf("No Slot Found with id %d", slot.SlotId)
	}
	return nil
}
func (r *SlotRepo) ListAvailableSlots() ([]domain.Slot, error) {
	rows, err := r.db.Query("SELECT slotid,slottype,isfree FROM slots WHERE isfree=true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []domain.Slot
	for rows.Next() {
		var s domain.Slot
		if err := rows.Scan(&s.SlotId, &s.SlotType, &s.IsFree); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, nil
}
func (r *SlotRepo) FindSlotByType(slottype string) ([]domain.Slot, error) {
	var Slots []domain.Slot
	rows, err := r.db.Query("SELECT slotid, slottype, isfree FROM slots WHERE slottype=? AND isfree=true", slottype)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slot domain.Slot
		if err := rows.Scan(&slot.SlotId, &slot.SlotType, &slot.IsFree); err != nil {
			return nil, err
		}
		Slots = append(Slots, slot)
	}
	return Slots, nil
}

func (r *SlotRepo) FindSlotTypebyID(SlotId int) string {
	var slottype string
	row := r.db.QueryRow("SELECT slottype from slots WHERE slotid=?", SlotId)
	err := row.Scan(&slottype)
	if err != nil {
		fmt.Println("Error fetching the Slottype by id", err)
		return " "
	}

	return slottype

}
func (r *SlotRepo) FindSlotByID(SlotId int) (*domain.Slot, error) {
	var Slot domain.Slot
	row := r.db.QueryRow("SELECT slotid, slottype, isfree FROM slots WHERE slotid = ?", SlotId)
	err := row.Scan(&Slot.SlotId, &Slot.SlotType, &Slot.IsFree)
	if err != nil {
		return nil, err
	}
	return &Slot, nil
}
