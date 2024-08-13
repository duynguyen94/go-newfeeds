package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	dbhost = "localhost"
	dbport = "3306"
	dbuser = "mysql"
	dbpass = "mysql"
	dbname = "newsfeed"
)

func initDBConn() (*sql.DB, error) {
	dataSource := dbuser + ":" + dbpass + "@" + dbhost + ":" + dbport + "/" + dbname
	db, err := sql.Open("mysql", dataSource)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}

func CreateNewUser(db *sql.DB, newUser UserRecord) (int64, error) {
	stmtIn, err := db.Prepare("INSERT INTO user (first_name, last_name, user_name, email, salt, hashed_password) VALUES (?, ?, ?, ?, ?, ?, ?)")
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

func GetUserRecord(db *sql.DB, id int) (UserRecord, error) {
	stmtOut, err := db.Prepare("SELECT first_name, last_name, user_name, email, salt, hashed_password FROM user WHERE id = ?")
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

func OverwriteUserRecord(db *sql.DB, user *UserRecord) error {
	stmtIn, err := db.Prepare("UPDATE `user` SET first_name = ?, last_name = ?, user_name = ?, email = ?, salt = ?, hashed_password = ? WHERE id = ? ")
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
		&user.Id,
	)

	return nil
}
