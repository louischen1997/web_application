package config

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

var err error

func SetupDatabaseConnection() {

	ps := os.Getenv("ps")
	ur := os.Getenv("ur")
	hs := os.Getenv("hs")
	dbn := os.Getenv("dbn")

	dsn := ur + ":" + ps + "@tcp(" + hs + ":3306)/" + dbn + "?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "root:chenYTCfor6225~!@tcp(localhost:3306)/Golangapi?charset=utf8mb4&parseTime=True&loc=Local"
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
