package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

//DB Declare database
var DB *gorm.DB

func init() {
	database, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		log.Fatal("Failed to connect to database!")
		panic("Check if test.db file exists!")
	}
	database.AutoMigrate(Account{}, Credential{})

	DB = database

}

//GetDB Returns instance of database
func GetDB() *gorm.DB {
	return DB
}
