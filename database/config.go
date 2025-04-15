package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"jwt-banking/models"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost port=5432 user=postgres password=pranish dbname=banking_system sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = db
	db.AutoMigrate(&models.User{})

	return db, nil
}
