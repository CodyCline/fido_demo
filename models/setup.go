package models

import (
	"github.com/jinzhu/gorm"
)

//DB Declare database
var DB *gorm.DB

func init() {
	database, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		panic("Check if test.db file exists!")
	}
	database.AutoMigrate(Account{}, Credential{})

	DB = database

}

//GetDB Returns instance of database
func GetDB() *gorm.DB {
	return DB
}
