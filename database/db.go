package database

import (
	"fmt"
	"os"
	"time"

	"geospatial--backend/models"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func InitDB() (*gorm.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost port=5432 user=postgres password=root dbname=postgres sslmode=disable"
	}

	// Open database connection
	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Enable logger
	db.LogMode(true)

	// Set connection pool settings
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)

	// Auto migrate the schema
	db.AutoMigrate(&models.User{}, &models.GeoData{})

	// Add any necessary indexes
	db.Model(&models.GeoData{}).AddIndex("idx_geodata_user_id", "user_id")

	// Log success message
	fmt.Println("Database connection established successfully")

	return db, nil
}
func CleanupDB(db *gorm.DB) error {
	// Hard delete all records (be careful with this!)
	err := db.Unscoped().Exec("DELETE FROM geo_data").Error
	if err != nil {
		return err
	}

	// Reset sequence if needed
	err = db.Exec("ALTER SEQUENCE geo_data_id_seq RESTART WITH 1").Error
	if err != nil {
		return err
	}

	return nil
}
