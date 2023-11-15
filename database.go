package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() *gorm.DB {
	//? Get the environment variables
	err := godotenv.Load()

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := "user=" + dbUsername + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"

	//? Conneect to db
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//? Create the database if it doesn't exist
	//? Check if the database exists
	rows, err := db.Raw("SELECT 1 FROM pg_database WHERE datname = ?", dbName).Rows()
	if err != nil {
		log.Fatalf("Failed to check if database exists: %v", err)
	}

	//? If the database doesn't exist, create it
	if !rows.Next() {
		err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
	}

	//? Enable the uuid-ossp extension
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	}

	//? Auto-migrate the User model
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate User model: %v", err)
	}

	return db
}
