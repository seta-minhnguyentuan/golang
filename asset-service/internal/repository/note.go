package repository

import (
	"asset-service/internal/models"

	"gorm.io/gorm"
)

type NoteRepository interface {
	CreateNote(note *models.Note) error
	ListNotes() ([]models.Note, error)
	ListNotesByUserAccess(userID string) ([]models.Note, error)
	GetNote(id string) (models.Note, error)
	UpdateNote(id string, note any) (models.Note, error)
	DeleteNote(id string) error
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

func (r *noteRepository) ListNotes() ([]models.Note, error) {
	var notes []models.Note
	if err := r.db.Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *noteRepository) ListNotesByUserAccess(userID string) ([]models.Note, error) {
	var notes []models.Note
	err := r.db.Preload("Sharings").
		Joins("JOIN folders ON notes.folder_id = folders.id").
		Where("folders.owner_id = ? OR folders.id IN (?)",
			userID,
			r.db.Table("folder_sharings").Select("folder_id").Where("user_id = ?", userID)).
		Find(&notes).Error
	return notes, err
}

func (r *noteRepository) GetNote(id string) (models.Note, error) {
	var note models.Note
	if err := r.db.First(&note, "id = ?", id).Error; err != nil {
		return models.Note{}, err
	}
	return note, nil
}

func (r *noteRepository) UpdateNote(id string, note any) (models.Note, error) {
	var updatedNote models.Note
	if err := r.db.Model(&updatedNote).Where("id = ?", id).Updates(note).Error; err != nil {
		return models.Note{}, err
	}
	return updatedNote, nil
}

func (r *noteRepository) DeleteNote(id string) error {
	return r.db.Delete(&models.Note{}, "id = ?", id).Error
}
