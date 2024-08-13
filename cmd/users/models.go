package main

import (
	"time"
)

const DOBLayout = "2006-01-02"
const saltSize = 16

type UserRecord struct {
	Id         int    `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	UserName   string `json:"userName"`
	Password   string `json:"password"`
	DOB        string `json:"dob"`
	DOBDate    time.Time
	salt       []byte
	hashedPass string
}

func (u *UserRecord) DOBtoDate() error {
	t, err := time.Parse(DOBLayout, u.DOB)
	if err != nil {
		return err
	}

	u.DOBDate = t
	return nil
}

func (u *UserRecord) HashPassword() {
	salt := genRandomSalt(saltSize)
	hashedPass := hashPassword(u.Password, salt)

	u.salt = salt
	u.hashedPass = hashedPass
}

func (u *UserRecord) IsMatchPassword(curPass string) bool {
	return isPassMatch(u.hashedPass, curPass, u.salt)
}
