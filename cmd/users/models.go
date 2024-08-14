package main

import (
	"time"
)

const DOBLayout = "2006-01-02"
const saltSize = 16

type UserRecord struct {
	Id         int    `json:"id,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Email      string `json:"email,omitempty"`
	UserName   string `json:"userName,omitempty"`
	Password   string `json:"password,omitempty"`
	DOB        string `json:"dob,omitempty"`
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

func (u *UserRecord) Merge(updateRecord *UserRecord) {
	if updateRecord.FirstName != "" {
		u.FirstName = updateRecord.FirstName
	}

	if updateRecord.LastName != "" {
		u.LastName = updateRecord.LastName
	}

	if updateRecord.Password != "" {
		u.Password = updateRecord.Password
		u.HashPassword()
	}

	if updateRecord.DOB != "" {
		u.DOB = updateRecord.DOB
	}

}

type PostRecord struct {
	Id               int    `json:"id"`
	ContentText      string `json:"text"`
	ContentImagePath string `json:"imagePath"`
	CreatedAt        string `json:"createdAt"`
	UserId           int    `json:"userId"`
}
