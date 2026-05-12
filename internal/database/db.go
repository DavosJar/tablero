package database

import (
	"fmt"
	"os"
	"tablero/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getPostgresDSN() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	if host == "" || user == "" || name == "" {
		return ""
	}
	if port == "" {
		port = "5432"
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, pass, name)
}

func InitDB() (*gorm.DB, error) {
	dsn := getPostgresDSN()
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL or DB_HOST/DB_USER/DB_PASSWORD/DB_NAME is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Column{}, &models.Task{})
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
		var count int64
		db.Model(&models.Column{}).Where("name = ?", col.Name).Count(&count)
		if count == 0 {
			db.Create(&col)
		}
	}

	return nil
}
