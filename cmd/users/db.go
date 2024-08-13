package main

import (
	"database/sql"
)

type UserDBModel struct {
	DB *sql.DB
}

func (m UserDBModel) CreateNewUser(newUser UserRecord) (int64, error) {
	stmtIn, err := m.DB.Prepare("INSERT INTO user (first_name, last_name, user_name, email, salt, hashed_password, dob) VALUES (?, ?, ?, ?, ?, ?, ?)")
	defer stmtIn.Close()

	if err != nil {
		return -1, err
	}

	result, err := stmtIn.Exec(
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.UserName,
		&newUser.Email,
		&newUser.salt,
		&newUser.hashedPass,
		&newUser.DOB,
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

func (m UserDBModel) GetUserRecord(id int) (UserRecord, error) {
	stmtOut, err := m.DB.Prepare("SELECT first_name, last_name, user_name, email, salt, hashed_password FROM user WHERE id = ?")
	defer stmtOut.Close()

	var user UserRecord

	if err != nil {
		return user, err
	}

	err = stmtOut.QueryRow(id).Scan(
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Email,
		&user.salt,
		&user.hashedPass,
	)

	return user, err
}

func (m UserDBModel) GetUserRecordByUsername(username string) (UserRecord, error) {
	stmtOut, err := m.DB.Prepare("SELECT first_name, last_name, user_name, email, salt, hashed_password FROM user WHERE user_name = ?")
	defer stmtOut.Close()

	var user UserRecord

	if err != nil {
		return user, err
	}

	err = stmtOut.QueryRow(username).Scan(
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Email,
		&user.salt,
		&user.hashedPass,
	)

	return user, err
}

func (m UserDBModel) OverwriteUserRecord(id int, user *UserRecord) error {
	stmtIn, err := m.DB.Prepare("UPDATE `user` SET first_name = ?, last_name = ?, user_name = ?, email = ?, salt = ?, hashed_password = ? WHERE id = ? ")
	defer stmtIn.Close()

	if err != nil {
		return err
	}

	_, err = stmtIn.Exec(
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Email,
		&user.salt,
		&user.hashedPass,
		id,
	)

	return nil
}
