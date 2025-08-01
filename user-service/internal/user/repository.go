package user

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FetchAll(ctx context.Context) ([]*User, error)
}

type GormRepository struct {
	DB *gorm.DB
}

func (r *GormRepository) Create(ctx context.Context, user *User) (*User, error) {
	if err := r.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormRepository) FetchAll(ctx context.Context) ([]*User, error) {
	var users []*User
	if err := r.DB.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
