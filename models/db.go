package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

// DB is the current database instance
var DB *gorm.DB

// SetupDB establishes a database connection and run auto-migrations
func SetupDB() *gorm.DB {
	conn := os.Getenv("DB_CONN")
	if conn == "" {
		conn = "host=localhost port=5432 user=postgres dbname=go_users_api_development password=postgres sslmode=disable"
	}

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		fmt.Printf("Error connecting to the database %s \n", err)
		panic("failed to connect database")
	}

	DB = db
	runDbMigrations()
	return db
}

func runDbMigrations() {
	DB.AutoMigrate(User{})
}
