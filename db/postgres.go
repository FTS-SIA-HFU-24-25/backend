package db

import (
	"log"
	"os"
	"sia/backend/tools"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Database *gorm.DB
)

func InitDB() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "[DATABASE]\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second * 3,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	Database, err := gorm.Open(postgres.Open(tools.POSTGRES_URI), &gorm.Config{
		TranslateError: true,
		Logger:         newLogger,
	})
	if err != nil {
		panic(err)
	}

	tools.Log("[DATABASE]", "Connected to PostgreSQL")
	// Database.AutoMigrate(&models.Hardware{}, &models.EcgData{}, &models.GpsData{}, &models.TemperatureData{})

	return Database
}
