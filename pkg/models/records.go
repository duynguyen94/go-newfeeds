package models

import (
	"github.com/duynguyen94/go-newfeeds/pkg/utils"
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
	salt := utils.GenRandomSalt(saltSize)
	hashedPass := utils.HashPassword(u.Password, salt)

	u.salt = salt
	u.hashedPass = hashedPass
}

func (u *UserRecord) IsMatchPassword(curPass string) bool {
	return utils.IsPassMatch(u.hashedPass, curPass, u.salt)
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
	Id               int    `json:"id,omitempty"`
	ContentText      string `json:"text,omitempty"`
	ContentImagePath string `json:"imagePath"`
	CreatedAt        string `json:"createdAt,omitempty"`
	UserId           int    `json:"userId,omitempty"`
	DownloadUrl      string `json:"downloadUrl,omitempty"`
}

func (p *PostRecord) GenSignedUrl(storage ImagePostStorageModel, expiration time.Duration) error {
	signedUrl, err := storage.GetSignedUrl(p.ContentImagePath, expiration)
	if err != nil {
		return err
	}

	p.DownloadUrl = signedUrl
	return nil
}

func (p *PostRecord) Merge(newPost *PostRecord) {
	if newPost.ContentText != "" {
		p.ContentText = newPost.ContentText
	}

	// What about image???
}

type CommentRecord struct {
	ContentText string `json:"content"`
	UserId      int    `json:"userId"`
}
