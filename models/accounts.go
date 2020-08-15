package models

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//Account A simple user account with authenticators
type Account struct {
	gorm.Model
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

//AddCredential Associates a credential with a user
func (a *Account) AddCredential(cred webauthn.Credential) {
	credJSON, err := json.Marshal(cred)
	if err != nil {
	}

	newCredential := Credential{
		Details:    credJSON,
		FKUsername: a.Username,
		AAGUID:     cred.Authenticator.AAGUID,
		SignCount:  cred.Authenticator.SignCount,
	}
	GetDB().Save(&newCredential)

}

//GetCredentials gets all the credentials associated with the user.
func GetCredentials(user string) []*Credential {
	credentials := make([]*Credential, 0)
	err := GetDB().Table("credentials").Where("fk_username = ?", user).Find(&credentials).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return credentials
}

//WebAuthnCredentials returns credentials owned by the user
func (a Account) WebAuthnCredentials() []webauthn.Credential {
	var creds []Credential
	credentialList := []webauthn.Credential{}
	GetDB().Table("credentials").Where("fk_username = ?", a.Username).Find(&creds)
	for _, cred := range creds {
		oneCred := webauthn.Credential{}
		json.Unmarshal(cred.Details, &oneCred)
		credentialList = append(credentialList, oneCred)
	}

	return credentialList
}

//CredentialExcludeList Returns array with users Credentials
func (a Account) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	var credentials []Credential

	GetDB().Where("fk_username = ?", a.Username).Find(&credentials)

	for _, cred := range credentials {
		oneCred := webauthn.Credential{}
		json.Unmarshal(cred.Details, &oneCred)
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: oneCred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}
	return credentialExcludeList
}

//WebAuthnID Get webauthn id
func (a Account) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(a.ID))
	return buf
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
