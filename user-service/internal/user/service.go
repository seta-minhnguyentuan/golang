package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repo Repository
}

func (s *Service) CreateUser(ctx context.Context, username, email, password, role string) (*User, error) {
	if role != "manager" && role != "member" {
		return nil, errors.New("invalid role")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user := &User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
	}
	return s.Repo.Create(ctx, user)
}

func (s *Service) Login(ctx context.Context, email, password string) (*User, error) {
	user, err := s.Repo.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.Repo.FindByEmail(ctx, email)
}

func (s *Service) FetchUsers(ctx context.Context) ([]*User, error) {
	return s.Repo.FetchAll(ctx)
}
