package models

import (
	"crypto/rand"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

type baseModel struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

//Account A simple user account with authenticators
type Account struct {
	ID          uint         `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   *time.Time   `json:"-" sql:"index"`
	Username    string       `json:"username"` //Username or email
	Name        string       `json:"name"`
	Credentials []Credential `json:"-"`
}

//NewUser Create a new user in the database
func NewUser(username string, name string) *Account {
	user := &Account{}
	user.Username = username
	user.Name = name
	GetDB().Create(user)
	return user
}

//GetUser Get a user from the database using their username
func GetUser(username string) *Account {
	account := &Account{}
	GetDB().Table("accounts").Where("username = ?", username).First(account)
	if account.Username == "" {
		return nil
	}
	return account
}

func randomID() string {
	buf := make([]byte, 8)
	rand.Read(buf)
	return string(buf[:])
}

//Webauthn Methods
