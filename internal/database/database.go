package database

import (
	"fmt"
	"log"
	"os"

	"github.com/money-advice/receipt-backend/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establishes database connection
func Connect() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Reduce logging for Cloud Run
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// Migrate runs database migrations only if needed
func Migrate() error {
	// Check if migration is needed by checking if users table exists and has correct structure
	if !needsMigration() {
		log.Println("Database schema is up to date, skipping migration")
		return nil
	}

	log.Println("Running database migration...")
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")
	return nil
}

// needsMigration checks if database migration is needed
func needsMigration() bool {
	// Check if users table exists
	if !DB.Migrator().HasTable(&models.User{}) {
		log.Println("Users table does not exist, migration needed")
		return true
	}

	// Check if all required columns exist
	requiredColumns := []string{"id", "created_at", "updated_at", "deleted_at", "google_id"}
	for _, column := range requiredColumns {
		if !DB.Migrator().HasColumn(&models.User{}, column) {
			log.Printf("Column %s does not exist, migration needed", column)
			return true
		}
	}

	// Check if google_id has unique index
	if !DB.Migrator().HasIndex(&models.User{}, "google_id") {
		log.Println("Google ID index does not exist, migration needed")
		return true
	}

	return false
}

// Close closes database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}