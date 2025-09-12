package mysql

import (
	"database/sql"
	"parkingSlotManagement/internals/core/domain"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) GetByUsername(username string) (*domain.Admin, error) {
	var admin domain.Admin
	query := "SELECT id, username, password FROM users WHERE username = ?"
	err := r.db.QueryRow(query, username).Scan(&admin.ID, &admin.Username, &admin.Password)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
