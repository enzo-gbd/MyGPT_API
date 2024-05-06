package main

import (
	"log"

	"github.com/enzo-gbd/GBA/configs"
	"github.com/enzo-gbd/GBA/internal/db"
	"github.com/enzo-gbd/GBA/internal/models"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	database := db.InitDB(&config)
	err = database.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
