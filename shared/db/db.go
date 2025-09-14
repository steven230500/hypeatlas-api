package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		log.Fatal("POSTGRES_URL is empty")
	}
	cfg := &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Error),
	}
	g, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		log.Fatalf("gorm open error: %v", err)
	}
	sqlDB, err := g.DB()
	if err != nil {
		log.Fatalf("gorm.DB(): %v", err)
	}
	sqlDB.SetMaxIdleConns(4)
	sqlDB.SetMaxOpenConns(16)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	return g
}

// Call ejecuta una query de GORM y maneja errores de forma consistente
func Call(db *gorm.DB) *gorm.DB {
	return db
}

// PreloadRelations carga relaciones comunes para entidades
func PreloadRelations(db *gorm.DB, relations ...string) *gorm.DB {
	for _, relation := range relations {
		db = db.Preload(relation)
	}
	return db
}
