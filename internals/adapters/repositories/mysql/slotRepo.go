package mysql

import (
	"database/sql"
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
		return Wrap("error inserting slot", err)
	}
	return nil

}

func (r *SlotRepo) UpdateSlot(slot *domain.Slot) error {
	res, err := r.db.Exec("UPDATE slots SET slottype=?, isfree=? WHERE slotid=?",
		slot.SlotType, slot.IsFree, slot.SlotId)
	if err != nil {
		return Wrap("error executing update slot query", err)
	}
	row, err := res.RowsAffected()
	if err != nil {
		return Wrap("error checking rows affected for slot update", err)
	}
	if row == 0 {
		return ErrSlotNotFoundByID
	}
	return nil
}
func (r *SlotRepo) ListAvailableSlots() ([]domain.Slot, error) {
	rows, err := r.db.Query("SELECT slotid,slottype,isfree FROM slots WHERE isfree=true")
	if err != nil {
		return nil, Wrap("error fetching slots :", err)
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
		return nil, Wrap("error fetching slot by type :", err)
	}
	defer rows.Close()

	for rows.Next() {
		var slot domain.Slot
		if err := rows.Scan(&slot.SlotId, &slot.SlotType, &slot.IsFree); err != nil {
			return nil, ErrSlotNotFound
		}
		Slots = append(Slots, slot)
	}
	return Slots, nil
}

func (r *SlotRepo) FindSlotTypebyID(SlotId int) (string, error) {
	var slottype string
	row := r.db.QueryRow("SELECT slottype from slots WHERE slotid=?", SlotId)
	err := row.Scan(&slottype)
	if err != nil {

		if err == sql.ErrNoRows {
			return "", ErrSlotNotFound
		}
		return "", Wrap("error fetching slot type by ID", err)
	}
	return slottype, nil

}
func (r *SlotRepo) FindSlotByID(SlotId int) (*domain.Slot, error) {
	var Slot domain.Slot
	row := r.db.QueryRow("SELECT slotid, slottype, isfree FROM slots WHERE slotid = ?", SlotId)
	err := row.Scan(&Slot.SlotId, &Slot.SlotType, &Slot.IsFree)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSlotNotFound
		}
		return nil, Wrap("error scanning slot by ID", err)
	}
	return &Slot, nil

}
