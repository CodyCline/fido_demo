package models

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

//Credential Represents the data from a key in database serialized form
type Credential struct {
	ID         uint       `json:"id" gorm:"primary_key"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"-" sql:"index"`
	AAGUID     []byte     `json:"-" gorm:"size:255"`
	Details    []byte     `json:"-" gorm:"size:2048"`
	SignCount  uint32     `json:"sign_count"`
	FKUsername string     `json:"-"`
}

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

//GetCredential obtains a single credential given its aaguid
func GetCredential(aaguid []byte) (Credential, error) {
	credential := Credential{}
	if result := GetDB().Where("aa_guid = ?", aaguid).First(&credential); result.Error != nil {
		return credential, errors.New("Credential does not exist")
	}
	return credential, nil
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

//UpdateCredential updates the
func UpdateCredential(aaguid []byte, counter uint32) {
	credential := &Credential{}
	GetDB().
		Model(&credential).
		Where("aa_guid = ?", aaguid).
		Updates(map[string]interface{}{
			"sign_count": counter,
			"updated_at": time.Now(),
		})
	return
}

//Below are all webAuthn methods needed by the library to implement the protocol

//WebAuthnCredentials returns credentials owned by the user
func (a Account) WebAuthnCredentials() []webauthn.Credential {
	var creds []Credential
	credentialList := []webauthn.Credential{}
	GetDB().
		Table("credentials").
		Where("fk_username = ?", a.Username).
		Find(&creds)
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
