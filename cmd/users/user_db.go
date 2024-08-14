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

func (m UserDBModel) FollowUser(id int, friendId int) error {
	stmtIn, err := m.DB.Prepare("INSERT INTO `user_user` (fk_user_id, fk_follower_id) VALUES (?, ?)")
	defer stmtIn.Close()

	if err != nil {
		return err
	}

	_, err = stmtIn.Exec(id, friendId)
	return err
}

func (m UserDBModel) UnFollowUser(id int, friendId int) error {
	stmtIn, err := m.DB.Prepare("DELETE FROM `user_user` WHERE fk_user_id=? AND fk_follower_id=?")
	defer stmtIn.Close()

	if err != nil {
		return err
	}

	_, err = stmtIn.Exec(id, friendId)
	return err
}

func (m UserDBModel) ViewFollowers(id int) ([]UserRecord, error) {
	var followers []UserRecord

	stmtOut, err := m.DB.Prepare("SELECT first_name, last_name, email, user_name FROM `user_user` u_u LEFT JOIN `user` u ON(u_u.fk_follower_id = u.id) WHERE u_u.fk_user_id = ?")
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
		var u UserRecord

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

func (m UserDBModel) ViewFriendPost(id int) ([]PostRecord, error) {
	var posts []PostRecord

	stmtOut, err := m.DB.Prepare("SELECT id, content_text, content_image_path, created_at FROM `user_user` u_u LEFT JOIN `post` p ON u_u.fk_follower_id = p.fk_user_id WHERE u_u.fk_user_id = ? AND visible = 1")
	defer stmtOut.Close()

	if err != nil {
		return posts, err
	}

	rows, err := stmtOut.Query(id)
	defer rows.Close()

	if err != nil {
		return posts, err
	}

	for rows.Next() {
		var p PostRecord

		err := rows.Scan(&p.Id, &p.ContentText, &p.ContentImagePath, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	err = rows.Err()
	if err != nil {
		return posts, err
	}

	return posts, nil
}
