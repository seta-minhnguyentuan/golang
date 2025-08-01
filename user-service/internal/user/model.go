package user

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username     string    `gorm:"not null"`
	Email        string    `gorm:"unique;not null"`
	Role         string    `gorm:"type:VARCHAR(10);not null"`
	PasswordHash string    `gorm:"not null"`
}
