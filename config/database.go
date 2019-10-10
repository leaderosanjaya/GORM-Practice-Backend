package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	// postgres driver for gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

// ConnectDB to connect to DB
func ConnectDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("[DB Load Env] %s\n", err)
		fmt.Printf("Attempting to load online environment...\n")
	}

	host := os.Getenv("PG_HOST")
	username := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("PG_DBNAME")
	port, _ := strconv.Atoi(os.Getenv("PG_PORT"))

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, pass, dbname)

	db, err := gorm.Open("postgres", dbInfo)
	if err != nil {
		fmt.Printf("[DB ConnectDB] %s", err)
		return nil, err
	}
	fmt.Printf("Established connection successfully to DB %s\n", dbname)
	return db, nil
}
