package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Account struct {
	gorm.Model
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	credentials []Credential
}

type Credential struct {
	gorm.Model
	AAGUID    []byte `gorm:"size:255"`
	Details   []byte `gorm:"size:2048"`
	SignCount uint32
	FKName    string
}
