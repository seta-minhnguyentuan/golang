package services

import (
	"asset-service/internal/models"
	"asset-service/internal/repository"

	"github.com/google/uuid"
)

type FolderService interface {
	Create(name string) (any, error)
	GetByID(id string) (any, error)
	List() ([]any, error)
	Delete(id string) error
}

type folderService struct {
	repo repository.FolderRepository
}

func NewFolderService(repo repository.FolderRepository) FolderService {
	return &folderService{repo: repo}
}

func (s *folderService) Create(name string) (any, error) {
	folder := &models.Folder{
		Name: name,
	}

	err := s.repo.Create(folder)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *folderService) GetByID(id string) (any, error) {
	folderID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	folder, err := s.repo.GetByID(folderID)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *folderService) List() ([]any, error) {
	folders, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	result := make([]any, len(folders))
	for i, folder := range folders {
		result[i] = folder
	}

	return result, nil
}

func (s *folderService) Delete(id string) error {
	folderID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(folderID)
}
