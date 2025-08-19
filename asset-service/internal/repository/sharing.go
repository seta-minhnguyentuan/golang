package repository

import (
	"asset-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SharingRepository interface {
	ShareFolder(sharing *models.FolderSharing) error
	GetFolderSharing(folderID, userID uuid.UUID) (*models.FolderSharing, error)
	ListFolderSharings(folderID uuid.UUID) ([]models.FolderSharing, error)
	RevokeFolderSharing(folderID, userID uuid.UUID) error

	ShareNote(sharing *models.NoteSharing) error
	GetNoteSharing(noteID, userID uuid.UUID) (*models.NoteSharing, error)
	ListNoteSharings(noteID uuid.UUID) ([]models.NoteSharing, error)
	RevokeNoteSharing(noteID, userID uuid.UUID) error
}

type sharingRepository struct {
	db *gorm.DB
}

func NewSharingRepository(db *gorm.DB) SharingRepository {
	return &sharingRepository{db: db}
}

// Folder sharing methods
func (r *sharingRepository) ShareFolder(sharing *models.FolderSharing) error {
	// Check if sharing already exists and update it, otherwise create new
	var existingSharing models.FolderSharing
	err := r.db.Where("folder_id = ? AND user_id = ?", sharing.FolderID, sharing.UserID).First(&existingSharing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new sharing
		return r.db.Create(sharing).Error
	} else if err != nil {
		return err
	}

	// Update existing sharing
	existingSharing.Permission = sharing.Permission
	return r.db.Save(&existingSharing).Error
}

func (r *sharingRepository) GetFolderSharing(folderID, userID uuid.UUID) (*models.FolderSharing, error) {
	var sharing models.FolderSharing
	err := r.db.Where("folder_id = ? AND user_id = ?", folderID, userID).First(&sharing).Error
	if err != nil {
		return nil, err
	}
	return &sharing, nil
}

func (r *sharingRepository) ListFolderSharings(folderID uuid.UUID) ([]models.FolderSharing, error) {
	var sharings []models.FolderSharing
	err := r.db.Where("folder_id = ?", folderID).Find(&sharings).Error
	return sharings, err
}

func (r *sharingRepository) RevokeFolderSharing(folderID, userID uuid.UUID) error {
	return r.db.Where("folder_id = ? AND user_id = ?", folderID, userID).Delete(&models.FolderSharing{}).Error
}

// Note sharing methods
func (r *sharingRepository) ShareNote(sharing *models.NoteSharing) error {
	// Check if sharing already exists and update it, otherwise create new
	var existingSharing models.NoteSharing
	err := r.db.Where("note_id = ? AND user_id = ?", sharing.NoteID, sharing.UserID).First(&existingSharing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new sharing
		return r.db.Create(sharing).Error
	} else if err != nil {
		return err
	}

	// Update existing sharing
	existingSharing.Permission = sharing.Permission
	return r.db.Save(&existingSharing).Error
}

func (r *sharingRepository) GetNoteSharing(noteID, userID uuid.UUID) (*models.NoteSharing, error) {
	var sharing models.NoteSharing
	err := r.db.Where("note_id = ? AND user_id = ?", noteID, userID).First(&sharing).Error
	if err != nil {
		return nil, err
	}
	return &sharing, nil
}

func (r *sharingRepository) ListNoteSharings(noteID uuid.UUID) ([]models.NoteSharing, error) {
	var sharings []models.NoteSharing
	err := r.db.Where("note_id = ?", noteID).Find(&sharings).Error
	return sharings, err
}

func (r *sharingRepository) RevokeNoteSharing(noteID, userID uuid.UUID) error {
	return r.db.Where("note_id = ? AND user_id = ?", noteID, userID).Delete(&models.NoteSharing{}).Error
}
