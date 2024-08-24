package database

import (
	"database/sql"
	"github.com/duynguyen94/go-newfeeds/internal/payloads"
)

type UserDB interface {
	// New add new user and return latest ID
	New(userPayload *payloads.UserPayload) (int64, error)

	// GetById get user by ID
	GetById(id int) (payloads.UserPayload, error)

	// GetByUsername get user by Username
	GetByUsername(username string) (payloads.UserPayload, error)

	// Edit overwrite user record given new record
	Edit(id int, payload *payloads.UserPayload) error

	// Follow user follow other users
	Follow(id int, friendId int) error

	// UnFollow other users
	UnFollow(id int, friendId int) error

	// GetFollowers list all followers
	GetFollowers(id int) ([]payloads.UserPayload, error)
}

func NewUserDB(db *sql.DB) UserDB {
	return &userDB{db: db}
}

type userDB struct {
	db *sql.DB
}

func (m *userDB) New(userPayload *payloads.UserPayload) (int64, error) {
	stmtIn, err := m.db.Prepare("INSERT INTO user (first_name, last_name, user_name, email, salt, hashed_password, dob) VALUES (?, ?, ?, ?, ?, ?, ?)")
	defer stmtIn.Close()

	if err != nil {
		return -1, err
	}

	result, err := stmtIn.Exec(
		&userPayload.FirstName,
		&userPayload.LastName,
		&userPayload.UserName,
		&userPayload.Email,
		&userPayload.Salt,
		&userPayload.HashedPass,
		&userPayload.DOB,
	)
	if err != nil {
		return -1, err
	}

	// Return new employee
	newUserId, err := result.LastInsertId()
	if err != nil {
		return newUserId, err
	}

	return newUserId, nil

}

func (m *userDB) GetById(id int) (payloads.UserPayload, error) {
	stmtOut, err := m.db.Prepare("SELECT first_name, last_name, user_name, email, salt, hashed_password FROM user WHERE id = ?")
	defer stmtOut.Close()

	var user payloads.UserPayload

	if err != nil {
		return user, err
	}

	err = stmtOut.QueryRow(id).Scan(
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Email,
		&user.Salt,
		&user.HashedPass,
	)

	return user, err
}

func (m *userDB) GetByUsername(username string) (payloads.UserPayload, error) {
	stmtOut, err := m.db.Prepare("SELECT first_name, last_name, user_name, email, salt, hashed_password FROM user WHERE user_name = ?")
	defer stmtOut.Close()

	var user payloads.UserPayload

	if err != nil {
		return user, err
	}

	err = stmtOut.QueryRow(username).Scan(
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Email,
		&user.Salt,
		&user.HashedPass,
	)

	return user, err
}

func (m *userDB) Edit(id int, user *payloads.UserPayload) error {
	stmtIn, err := m.db.Prepare("UPDATE `user` SET first_name = ?, last_name = ?, user_name = ?, email = ?, salt = ?, hashed_password = ? WHERE id = ? ")
	defer stmtIn.Close()

	if err != nil {
		return err
	}

	_, err = stmtIn.Exec(
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Email,
		&user.Salt,
		&user.HashedPass,
		id,
	)

	return nil
}

func (m *userDB) Follow(id int, friendId int) error {
	stmtIn, err := m.db.Prepare("INSERT INTO `user_user` (fk_user_id, fk_follower_id) VALUES (?, ?)")
	defer stmtIn.Close()

	if err != nil {
		return err
	}

	_, err = stmtIn.Exec(id, friendId)
	return err
}

func (m *userDB) UnFollow(id int, friendId int) error {
	stmtIn, err := m.db.Prepare("DELETE FROM `user_user` WHERE fk_user_id=? AND fk_follower_id=?")
	defer stmtIn.Close()

	if err != nil {
		return err
	}

	_, err = stmtIn.Exec(id, friendId)
	return err
}

func (m *userDB) GetFollowers(id int) ([]payloads.UserPayload, error) {
	var followers []payloads.UserPayload

	stmtOut, err := m.db.Prepare("SELECT first_name, last_name, email, user_name FROM `user_user` u_u LEFT JOIN `user` u ON(u_u.fk_follower_id = u.id) WHERE u_u.fk_user_id = ?")
	defer stmtOut.Close()

	if err != nil {
		return followers, err
	}

	rows, err := stmtOut.Query(id)
	defer rows.Close()

	if err != nil {
		return followers, err
	}

	for rows.Next() {
		var u payloads.UserPayload

		err := rows.Scan(&u.FirstName, &u.LastName, &u.Email, &u.UserName)
		if err != nil {
			return nil, err
		}

		followers = append(followers, u)
	}

	err = rows.Err()
	if err != nil {
		return followers, err
	}

	return followers, nil
}
