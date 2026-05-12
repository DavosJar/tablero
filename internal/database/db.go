package database

import (
	"os"
	"tablero/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")

	var db *gorm.DB
	var err error

	if databaseURL != "" {
		// PostgreSQL
		db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	} else {
		// SQLite
		db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}

	// Auto-migrate
	db.AutoMigrate(&models.Column{}, &models.Task{})

	// Crear columnas por defecto
	createDefaultColumns(db)

	return db, nil
}

func createDefaultColumns(db *gorm.DB) error {
	defaultColumns := []models.Column{
		{Name: "Por hacer", Order: 1},
		{Name: "En progreso", Order: 2},
		{Name: "Hecho", Order: 3},
	}

	for _, col := range defaultColumns {
		// Check if column already exists
		var count int64
		db.Model(&models.Column{}).Where("name = ?", col.Name).Count(&count)
		if count == 0 {
			db.Create(&col)
		}
	}

	return nil
}
