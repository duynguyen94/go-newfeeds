package main

import (
	"time"
)

const DOBLayout = "2006-01-02"

type UserRecord struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	UserName  string `json:"userName"`
	Password  string `json:"password"`
	DOB       string `json:"dob"`
	DOBDate   time.Time
}

func (u *UserRecord) DOBtoDate() error {
	t, err := time.Parse(DOBLayout, u.DOB)
	if err != nil {
		return err
	}

	u.DOBDate = t
	return nil
}
