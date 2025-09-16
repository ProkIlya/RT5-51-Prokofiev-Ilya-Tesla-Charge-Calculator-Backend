package main

import (
	"log"
	"tesla-app/internal/app/ds"
	"tesla-app/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// Миграция схемы
	err = db.AutoMigrate(
		&ds.User{},
		&ds.DrivingScenario{},
		&ds.TripApplication{},
		&ds.TripScenario{},
	)
	if err != nil {
		log.Fatal("cant migrate db: ", err)
	}
	log.Println("Migration completed successfully")

}
