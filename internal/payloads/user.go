package payloads

import (
	"github.com/duynguyen94/go-newfeeds/internal/utils"
	"time"
)

const DOBLayout = "2006-01-02"
const saltSize = 16

type UserPayload struct {
	Id         int    `json:"id,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Email      string `json:"email,omitempty"`
	UserName   string `json:"userName,omitempty"`
	Password   string `json:"password,omitempty"`
	DOB        string `json:"dob,omitempty"`
	DOBDate    time.Time
	Salt       []byte
	HashedPass string
}

func (u *UserPayload) DOBtoDate() error {
	t, err := time.Parse(DOBLayout, u.DOB)
	if err != nil {
		return err
	}

	u.DOBDate = t
	return nil
}

func (u *UserPayload) HashPassword() {
	salt := utils.GenRandomSalt(saltSize)
	hashedPass := utils.HashPassword(u.Password, salt)

	u.Salt = salt
	u.HashedPass = hashedPass
}

func (u *UserPayload) IsMatchPassword(curPass string) bool {
	return utils.IsPassMatch(u.HashedPass, curPass, u.Salt)
}

func (u *UserPayload) Merge(updateRecord *UserPayload) {
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
