package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

var err error

func SetupDatabaseConnection() {
	// errENV := godotenv.Load()
	// if errENV != nil {
	// 	panic("FAILED to LOAD env file")
	// }

	dsn := "newur:ytc6225forclass@tcp(localhost:3306)/Goapi?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to close connection")
	}
	dbSQL.Close()
}
