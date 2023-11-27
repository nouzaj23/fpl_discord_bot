package database

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

var database *gorm.DB

func InitDB() *gorm.DB {
	godotenv.Load()
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_NAME")), &gorm.Config{})
	if err != nil {
		log.Fatalf("%v %v Failed to open database: %v", log.Ldate, log.Ltime, err)
	}
	database = db
	return db
}

func GetDB() *gorm.DB {
	return database
}
