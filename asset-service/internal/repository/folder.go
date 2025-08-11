package repository

import (
	"asset-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FolderRepository interface {
	Create(folder *models.Folder) error
	GetByID(id uuid.UUID) (*models.Folder, error)
	List() ([]models.Folder, error)
	Delete(id uuid.UUID) error
}

type folderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) FolderRepository {
	return &folderRepository{db: db}
}

func (r *folderRepository) Create(folder *models.Folder) error {
	return r.db.Create(folder).Error
}

func (r *folderRepository) GetByID(id uuid.UUID) (*models.Folder, error) {
	var folder models.Folder
	err := r.db.Where("id = ?", id).First(&folder).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *folderRepository) List() ([]models.Folder, error) {
	var folders []models.Folder
	err := r.db.Find(&folders).Error
	return folders, err
}

func (r *folderRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Folder{}, id).Error
}
