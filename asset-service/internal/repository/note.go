package repository

import (
	"asset-service/internal/models"

	"gorm.io/gorm"
)

type NoteRepository interface {
	CreateNote(note *models.Note) error
}

type noteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &noteRepository{db: db}
}

func (r *noteRepository) CreateNote(note *models.Note) error {
	return r.db.Create(note).Error
}
