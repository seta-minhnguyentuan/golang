package services

import (
	"asset-service/internal/models"
	"asset-service/internal/repository"
	"fmt"

	"github.com/google/uuid"
)

type NoteService interface {
	CreateNote(name string, content string, folderId uuid.UUID, userID uuid.UUID) (*models.Note, error)
	GetNote(id string, userID uuid.UUID) (models.Note, error)
	ListNotes(userID uuid.UUID) ([]models.Note, error)
	UpdateNote(id string, note any, userID uuid.UUID) (models.Note, error)
	DeleteNote(id string, userID uuid.UUID) error
}

type noteService struct {
	repo       repository.NoteRepository
	folderRepo repository.FolderRepository
}

func NewNoteService(noteRepo repository.NoteRepository, folderRepo repository.FolderRepository) NoteService {
	return &noteService{repo: noteRepo, folderRepo: folderRepo}
}

func (s *noteService) CreateNote(name string, content string, folderId uuid.UUID, userID uuid.UUID) (*models.Note, error) {
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

func (s *noteService) ListNotes(userID uuid.UUID) ([]models.Note, error) {
	notes, err := s.repo.ListNotesByUserAccess(userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	return notes, nil
}

func (s *noteService) GetNote(id string, userID uuid.UUID) (models.Note, error) {
	note, err := s.repo.GetNote(id)
	if err != nil {
		return models.Note{}, fmt.Errorf("failed to get note: %w", err)
	}

	// Check if user has access to the folder containing this note
	folder, err := s.folderRepo.GetFolderByID(note.FolderID)
	if err != nil {
		return models.Note{}, fmt.Errorf("failed to verify folder access: %w", err)
	}

	// Check if user owns the folder or has it shared with them
	hasAccess := folder.OwnerID == userID
	if !hasAccess {
		for _, sharing := range folder.Sharings {
			if sharing.UserID == userID {
				hasAccess = true
				break
			}
		}
	}

	if !hasAccess {
		return models.Note{}, fmt.Errorf("access denied: you don't have permission to view this note")
	}

	return note, nil
}

func (s *noteService) UpdateNote(id string, note any, userID uuid.UUID) (models.Note, error) {
	// First get the existing note to check folder access
	existingNote, err := s.repo.GetNote(id)
	if err != nil {
		return models.Note{}, fmt.Errorf("failed to get note: %w", err)
	}

	// Check if user has write access to the folder containing this note
	folder, err := s.folderRepo.GetFolderByID(existingNote.FolderID)
	if err != nil {
		return models.Note{}, fmt.Errorf("failed to verify folder access: %w", err)
	}

	// Check if user owns the folder or has write permission shared with them
	hasWriteAccess := folder.OwnerID == userID
	if !hasWriteAccess {
		for _, sharing := range folder.Sharings {
			if sharing.UserID == userID && sharing.Permission == "write" {
				hasWriteAccess = true
				break
			}
		}
	}

	if !hasWriteAccess {
		return models.Note{}, fmt.Errorf("access denied: you don't have write permission for this note")
	}

	updatedNote, err := s.repo.UpdateNote(id, note)
	if err != nil {
		return models.Note{}, fmt.Errorf("failed to update note: %w", err)
	}

	return updatedNote, nil
}

func (s *noteService) DeleteNote(id string, userID uuid.UUID) error {
	// First get the existing note to check folder access
	existingNote, err := s.repo.GetNote(id)
	if err != nil {
		return fmt.Errorf("failed to get note: %w", err)
	}

	// Check if user has write access to the folder containing this note
	folder, err := s.folderRepo.GetFolderByID(existingNote.FolderID)
	if err != nil {
		return fmt.Errorf("failed to verify folder access: %w", err)
	}

	// Check if user owns the folder or has write permission shared with them
	hasWriteAccess := folder.OwnerID == userID
	if !hasWriteAccess {
		for _, sharing := range folder.Sharings {
			if sharing.UserID == userID && sharing.Permission == "write" {
				hasWriteAccess = true
				break
			}
		}
	}

	if !hasWriteAccess {
		return fmt.Errorf("access denied: you don't have write permission to delete this note")
	}

	if err := s.repo.DeleteNote(id); err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	return nil
}
