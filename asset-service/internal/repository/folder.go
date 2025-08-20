package repository

import (
	"asset-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FolderRepository interface {
	CreateFolder(folder *models.Folder) error
	GetFolderByID(id uuid.UUID) (*models.Folder, error)
	ListFolders() ([]models.Folder, error)
	ListFoldersByOwner(ownerID uuid.UUID) ([]models.Folder, error)
	DeleteFolder(id uuid.UUID) error
}

type folderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) FolderRepository {
	return &folderRepository{db: db}
}

func (r *folderRepository) CreateFolder(folder *models.Folder) error {
	return r.db.Create(folder).Error
}

func (r *folderRepository) GetFolderByID(id uuid.UUID) (*models.Folder, error) {
	var folder models.Folder
	err := r.db.Preload("Notes").Preload("Sharings").Where("id = ?", id).First(&folder).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *folderRepository) ListFolders() ([]models.Folder, error) {
	var folders []models.Folder
	err := r.db.Preload("Notes").Preload("Sharings").Find(&folders).Error
	return folders, err
}

func (r *folderRepository) ListFoldersByOwner(ownerID uuid.UUID) ([]models.Folder, error) {
	var folders []models.Folder
	err := r.db.Preload("Notes").Preload("Sharings").Where("owner_id = ?", ownerID).Find(&folders).Error
	return folders, err
}

func (r *folderRepository) DeleteFolder(id uuid.UUID) error {
	return r.db.Delete(&models.Folder{}, id).Error
}
