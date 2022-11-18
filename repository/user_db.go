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

func (u userRepositoryDB) CreateUser(email string, password string, secret string) (*User, error) {
	// Insert document into database
	insert, err := u.db.Exec("INSERT INTO users (email, password, secret) VALUES (?, ?, ?)", email, password, secret)
	if err != nil {
		return nil, err
	}
	userId, err := insert.LastInsertId()

	// Create user object
	var user = User{
		Id:       userId,
		Email:    email,
		Password: password,
		Secret:   secret,
	}
	return &user, nil
}
func (u userRepositoryDB) CheckUser(email string) (*User, error) {
	result, err := u.db.Query("SELECT id, email, password, secret FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	var user User
	for result.Next() {
		if err := result.Scan(&user.Id, &user.Email, &user.Password, &user.Secret); err != nil {
			return nil, err
		}
	}
	defer result.Close()

	return &user, nil
}

func (u userRepositoryDB) GetUser(id int64) (*User, error) {
	result, err := u.db.Query("SELECT id, email, password, secret FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	var user User
	for result.Next() {
		if err := result.Scan(&user.Id, &user.Email, &user.Password, &user.Secret); err != nil {
			return nil, err
		}
	}

	defer result.Close()
	return &user, nil
}
