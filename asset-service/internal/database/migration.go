package database

import (
	"asset-service/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return err
	}

	return db.AutoMigrate(
		&models.Folder{},
		&models.Note{},
		&models.FolderSharing{},
		&models.NoteSharing{},
	)
}
