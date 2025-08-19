package services

import (
	"asset-service/internal/models"
	"asset-service/internal/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type SharingService interface {
	ShareFolder(folderID, userID uuid.UUID, permission models.Permission, ownerID uuid.UUID) error
	RevokeFolderSharing(folderID, userID uuid.UUID, ownerID uuid.UUID) error
	GetFolderSharing(folderID, userID uuid.UUID) (*models.FolderSharing, error)
	ListFolderSharings(folderID uuid.UUID, ownerID uuid.UUID) ([]models.FolderSharing, error)

	ShareNote(noteID, userID uuid.UUID, permission models.Permission, ownerID uuid.UUID) error
	RevokeNoteSharing(noteID, userID uuid.UUID, ownerID uuid.UUID) error
	GetNoteSharing(noteID, userID uuid.UUID) (*models.NoteSharing, error)
	ListNoteSharings(noteID uuid.UUID, ownerID uuid.UUID) ([]models.NoteSharing, error)
}

type sharingService struct {
	sharingRepo repository.SharingRepository
	folderRepo  repository.FolderRepository
	noteRepo    repository.NoteRepository
}

func NewSharingService(sharingRepo repository.SharingRepository, folderRepo repository.FolderRepository, noteRepo repository.NoteRepository) SharingService {
	return &sharingService{
		sharingRepo: sharingRepo,
		folderRepo:  folderRepo,
		noteRepo:    noteRepo,
	}
}

// Folder sharing methods
func (s *sharingService) ShareFolder(folderID, userID uuid.UUID, permission models.Permission, ownerID uuid.UUID) error {
	// Verify folder exists and user is the owner
	folder, err := s.folderRepo.GetFolderByID(folderID)
	if err != nil {
		return fmt.Errorf("folder not found: %w", err)
	}

	if folder.OwnerID != ownerID {
		return errors.New("only the folder owner can share the folder")
	}

	// Don't allow owner to share with themselves
	if folder.OwnerID == userID {
		return errors.New("cannot share folder with yourself")
	}

	// Validate permission
	if permission != models.PermissionRead && permission != models.PermissionWrite {
		return errors.New("invalid permission type")
	}

	sharing := &models.FolderSharing{
		FolderID:   folderID,
		UserID:     userID,
		Permission: permission,
	}

	return s.sharingRepo.ShareFolder(sharing)
}

func (s *sharingService) RevokeFolderSharing(folderID, userID uuid.UUID, ownerID uuid.UUID) error {
	// Verify folder exists and user is the owner
	folder, err := s.folderRepo.GetFolderByID(folderID)
	if err != nil {
		return fmt.Errorf("folder not found: %w", err)
	}

	if folder.OwnerID != ownerID {
		return errors.New("only the folder owner can revoke folder sharing")
	}

	return s.sharingRepo.RevokeFolderSharing(folderID, userID)
}

func (s *sharingService) GetFolderSharing(folderID, userID uuid.UUID) (*models.FolderSharing, error) {
	return s.sharingRepo.GetFolderSharing(folderID, userID)
}

func (s *sharingService) ListFolderSharings(folderID uuid.UUID, ownerID uuid.UUID) ([]models.FolderSharing, error) {
	// Verify folder exists and user is the owner
	folder, err := s.folderRepo.GetFolderByID(folderID)
	if err != nil {
		return nil, fmt.Errorf("folder not found: %w", err)
	}

	if folder.OwnerID != ownerID {
		return nil, errors.New("only the folder owner can view folder sharings")
	}

	return s.sharingRepo.ListFolderSharings(folderID)
}

// Note sharing methods
func (s *sharingService) ShareNote(noteID, userID uuid.UUID, permission models.Permission, ownerID uuid.UUID) error {
	// Verify note exists and get the folder to check ownership
	note, err := s.noteRepo.GetNote(noteID.String())
	if err != nil {
		return fmt.Errorf("note not found: %w", err)
	}

	// Get folder to check ownership
	folder, err := s.folderRepo.GetFolderByID(note.FolderID)
	if err != nil {
		return fmt.Errorf("parent folder not found: %w", err)
	}

	if folder.OwnerID != ownerID {
		return errors.New("only the note owner can share the note")
	}

	// Don't allow owner to share with themselves
	if folder.OwnerID == userID {
		return errors.New("cannot share note with yourself")
	}

	// Validate permission
	if permission != models.PermissionRead && permission != models.PermissionWrite {
		return errors.New("invalid permission type")
	}

	sharing := &models.NoteSharing{
		NoteID:     noteID,
		UserID:     userID,
		Permission: permission,
	}

	return s.sharingRepo.ShareNote(sharing)
}

func (s *sharingService) RevokeNoteSharing(noteID, userID uuid.UUID, ownerID uuid.UUID) error {
	// Verify note exists and get the folder to check ownership
	note, err := s.noteRepo.GetNote(noteID.String())
	if err != nil {
		return fmt.Errorf("note not found: %w", err)
	}

	// Get folder to check ownership
	folder, err := s.folderRepo.GetFolderByID(note.FolderID)
	if err != nil {
		return fmt.Errorf("parent folder not found: %w", err)
	}

	if folder.OwnerID != ownerID {
		return errors.New("only the note owner can revoke note sharing")
	}

	return s.sharingRepo.RevokeNoteSharing(noteID, userID)
}

func (s *sharingService) GetNoteSharing(noteID, userID uuid.UUID) (*models.NoteSharing, error) {
	return s.sharingRepo.GetNoteSharing(noteID, userID)
}

func (s *sharingService) ListNoteSharings(noteID uuid.UUID, ownerID uuid.UUID) ([]models.NoteSharing, error) {
	// Verify note exists and get the folder to check ownership
	note, err := s.noteRepo.GetNote(noteID.String())
	if err != nil {
		return nil, fmt.Errorf("note not found: %w", err)
	}

	// Get folder to check ownership
	folder, err := s.folderRepo.GetFolderByID(note.FolderID)
	if err != nil {
		return nil, fmt.Errorf("parent folder not found: %w", err)
	}

	if folder.OwnerID != ownerID {
		return nil, errors.New("only the note owner can view note sharings")
	}

	return s.sharingRepo.ListNoteSharings(noteID)
}
