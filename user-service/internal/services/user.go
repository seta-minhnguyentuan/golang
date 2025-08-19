package services

import (
	"context"
	"errors"
	"user-service/internal/models"
	"user-service/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, username, email, password, role string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	FetchUsers(ctx context.Context) ([]*models.User, error)
}

type UserServiceImpl struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{Repo: repo}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, username, email, password, role string) (*models.User, error) {
	if role != "manager" && role != "member" {
		return nil, errors.New("invalid role")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user := &models.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
	}
	return s.Repo.Create(ctx, user)
}

func (s *UserServiceImpl) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.Repo.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.Repo.FindByEmail(ctx, email)
}

func (s *UserServiceImpl) FetchUsers(ctx context.Context) ([]*models.User, error) {
	return s.Repo.FetchAll(ctx)
}
