package database

import (
	"activity-tracker-api/models"
	"activity-tracker-api/models/activity"
	"activity-tracker-api/models/gym"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

func ConnectPostgresDb() *gorm.DB {
	connStr := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("Failed to connect to database.\n", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&activity.DbActivity{},
		&activity.DbGroupActivity{},
		&gym.GymAccount{},
		&gym.GymEquipment{},
	); err != nil {
		return nil
	}

	return db
}
