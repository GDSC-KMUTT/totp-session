package repository

import (
	"database/sql"
)

type userRepositoryDB struct {
	db *sql.DB
}

func NewRepositoryDB(db *sql.DB) userRepositoryDB {
	return userRepositoryDB{db: db}
}

func (u userRepositoryDB) CreateUser(email string, password string, secret string) (*User, *string, error) {
	return nil, nil, nil
}
func (u userRepositoryDB) CheckUser(email string) (*User, error) {
	return nil, nil
}
