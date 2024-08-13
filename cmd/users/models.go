package main

import (
	"time"
)

const DOBLayout = "2006-01-02"
const saltSize = 16

type UserRecord struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	UserName   string `json:"userName"`
	Password   string `json:"password"`
	DOB        string `json:"dob"`
	DOBDate    time.Time
	salt       string
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

	u.salt = string(salt[:])
	u.hashedPass = hashedPass
}