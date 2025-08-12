package services

import (
	"asset-service/internal/models"
	"asset-service/internal/repository"

	"github.com/google/uuid"
)

type FolderService interface {
	CreateFolder(name string) (any, error)
	GetFolderByID(id string) (any, error)
	ListFolders() ([]any, error)
	DeleteFolder(id string) error
}

type folderService struct {
	repo repository.FolderRepository
}

func NewFolderService(repo repository.FolderRepository) FolderService {
	return &folderService{repo: repo}
}

func (s *folderService) CreateFolder(name string) (any, error) {
	folder := &models.Folder{
		Name: name,
	}

	err := s.repo.CreateFolder(folder)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *folderService) GetFolderByID(id string) (any, error) {
	folderID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	folder, err := s.repo.GetFolderByID(folderID)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *folderService) ListFolders() ([]any, error) {
	folders, err := s.repo.ListFolders()
	if err != nil {
		return nil, err
	}

	result := make([]any, len(folders))
	for i, folder := range folders {
		result[i] = folder
	}

	return result, nil
}

func (s *folderService) DeleteFolder(id string) error {
	folderID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return s.repo.DeleteFolder(folderID)
}
