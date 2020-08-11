package models

import (
	"crypto/rand"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//Account A simple user account with authenticators
type Account struct {
	gorm.Model
	ID          string
	Username    string `json:"username"` //Username or email
	Name        string `json:"name"`
	Credentials []Credential
}

//Credential Represents the data from a key in database serialized form
type Credential struct {
	gorm.Model
	AAGUID     []byte `gorm:"size:255"`
	Details    []byte `gorm:"size:2048"`
	SignCount  uint32
	FKUsername string
}

//NewUser Create a new user in the database
func NewUser(username string, name string) *Account {
	user := &Account{}
	user.ID = randomID()
	user.Username = name
	user.Name = username
	GetDB().Create(user)
	return user
}

//GetUser Get a user from the database using their username
func GetUser(name string) (*Account, error) {
	account := &Account{}
	GetDB().Table("accounts").Where("name = ?", name).First(account)
	if account.Username == "" {
		return account, fmt.Errorf("error getting user '%s': does not exist", name)
	}
	return account, nil
}

func randomID() string {
	buf := make([]byte, 8)
	rand.Read(buf)
	return string(buf[:])
}

//Webauthn Methods

//WebAuthnID Get webauthn id
func (a Account) WebAuthnID() []byte {
	return []byte(a.ID)
}

//WebAuthnName Return the users username such as email address
func (a Account) WebAuthnName() string {
	return a.Username
}

//WebAuthnDisplayName Return the user's displayed name
func (a Account) WebAuthnDisplayName() string {
	return a.Name
}

// WebAuthnIcon is not (yet) implemented
func (a Account) WebAuthnIcon() string {
	return ""
}
