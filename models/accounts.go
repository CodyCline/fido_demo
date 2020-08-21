package models

import (
	"crypto/rand"
	"fmt"
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
func GetUser(username string) (*Account, error) {
	account := &Account{}
	GetDB().Table("accounts").Where("username = ?", username).First(account)
	if account.Username == "" {
		return account, fmt.Errorf("error getting user '%s': does not exist", username)
	}
	return account, nil
}

func randomID() string {
	buf := make([]byte, 8)
	rand.Read(buf)
	return string(buf[:])
}

//Webauthn Methods
