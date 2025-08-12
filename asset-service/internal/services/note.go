package services

import (
	"asset-service/internal/models"
	"asset-service/internal/repository"
	"fmt"

	"github.com/google/uuid"
)

type NoteService interface {
	CreateNote(name string, content string, folderId uuid.UUID) (*models.Note, error)
	// GetNote(id string) (models.Note, error)
	// ListNotes() ([]models.Note, error)
	// UpdateNote(id string, note any) (models.Note, error)
	// DeleteNote(id string) error
}

type noteService struct {
	repo repository.NoteRepository
}

func NewNoteService(noteRepo repository.NoteRepository) NoteService {
	return &noteService{repo: noteRepo}
}

func (s *noteService) CreateNote(name string, content string, folderId uuid.UUID) (*models.Note, error) {
	if name == "" {
		return nil, fmt.Errorf("note name cannot be empty")
	}

	note := &models.Note{
		Name:     name,
		Content:  content,
		FolderID: folderId,
	}

	if err := s.repo.CreateNote(note); err != nil {
		return nil, err
	}

	return note, nil
}
