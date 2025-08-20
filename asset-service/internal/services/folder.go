package services

import (
	"asset-service/internal/models"
	"asset-service/internal/repository"
	"fmt"

	"github.com/google/uuid"
)

type FolderService interface {
	CreateFolder(name string, userID uuid.UUID) (any, error)
	GetFolderByID(id string, userID uuid.UUID) (any, error)
	ListFolders(userID uuid.UUID) ([]any, error)
	DeleteFolder(id string, userID uuid.UUID) error
}

type folderService struct {
	repo repository.FolderRepository
}

func NewFolderService(repo repository.FolderRepository) FolderService {
	return &folderService{repo: repo}
}

func (s *folderService) CreateFolder(name string, userID uuid.UUID) (any, error) {
	folder := &models.Folder{
		Name:      name,
		OwnerID:   userID,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	err := s.repo.CreateFolder(folder)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *folderService) GetFolderByID(id string, userID uuid.UUID) (any, error) {
	folderID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	folder, err := s.repo.GetFolderByID(folderID)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("access denied: you don't have permission to view this folder")
	}

	return folder, nil
}

func (s *folderService) ListFolders(userID uuid.UUID) ([]any, error) {
	folders, err := s.repo.ListFoldersByOwnerOrShared(userID)
	if err != nil {
		return nil, err
	}

	result := make([]any, len(folders))
	for i, folder := range folders {
		result[i] = folder
	}

	return result, nil
}

func (s *folderService) DeleteFolder(id string, userID uuid.UUID) error {
	folderID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	// First check if user owns the folder
	folder, err := s.repo.GetFolderByID(folderID)
	if err != nil {
		return err
	}

	if folder.OwnerID != userID {
		return fmt.Errorf("only the folder owner can delete this folder")
	}

	return s.repo.DeleteFolder(folderID)
}
