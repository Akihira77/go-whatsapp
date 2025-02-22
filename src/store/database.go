package store

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/Akihira77/go_whatsapp/src/types"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Store struct {
	DB *gorm.DB
}

func NewStore() *Store {
	dbPath := os.Getenv("DB_SQLITE")
	slog.Info("database",
		"path", dbPath,
	)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	db.Exec("PRAGMA foreign_keys = ON;")

	sqlDB, err := db.DB()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	slog.Info("Database",
		"stats", sqlDB.Stats(),
	)

	return &Store{
		DB: db,
	}
}

func (s *Store) Migrate() {
	slog.Info("Perfoming db migrations...")
	err := s.DB.AutoMigrate(
		&types.User{},
		&types.UserContact{},
		&types.File{},
		&types.Message{},
		&types.Group{},
		&types.UserGroup{},
	)

	if err != nil {
		log.Fatalf("Error perform db migrations: %v", err)
	}

	slog.Info("DB migrations done...")
}
