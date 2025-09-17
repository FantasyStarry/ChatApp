package config

import (
	"chatapp/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	if GlobalConfig == nil {
		log.Fatal("Configuration not loaded. Please call LoadConfig first.")
	}

	dsn := GlobalConfig.GetDatabaseDSN()

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}

	// Configure connection pool based on config
	sqlDB.SetMaxIdleConns(GlobalConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(GlobalConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(GlobalConfig.Database.ConnMaxLifetime)

	DB = database
	log.Println("Database connected successfully")
}

func MigrateDatabase() {
	if DB == nil {
		log.Fatal("Database not connected. Please call ConnectDatabase first.")
	}

	err := DB.AutoMigrate(&models.User{}, &models.ChatRoom{}, &models.Message{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed")
}
