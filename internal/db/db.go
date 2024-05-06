// Package db is responsible for setting up and managing the database connection using GORM.
package db

import (
	"fmt"
	"log"

	"github.com/enzo-gbd/GBA/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes and returns a new GORM database object configured
// with the provided settings from the config parameter. It sets up the database
// connection using the Postgres driver. The function will terminate the program
// with an error if the database connection cannot be established.
func InitDB(config *configs.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%v",
		config.DBHost,
		config.DBUserName,
		config.DBUserPassword,
		config.DBName,
		config.DBPort,
		config.DBSSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("could not connect with the database: %v", err)
	}

	return db
}
