package models

import (
	"time"

	"github.com/google/uuid"
)

type Permission string

const (
	PermissionRead  Permission = "read"
	PermissionWrite Permission = "write"
)

type FolderSharing struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null" json:"userId"`
	Permission Permission `gorm:"type:varchar(16);not null;check:permission IN ('read','write')" json:"permission"`
	FolderID   uuid.UUID  `gorm:"type:uuid" json:"folderId"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

type NoteSharing struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null" json:"userId"`
	Permission Permission `gorm:"type:varchar(16);not null;check:permission IN ('read','write')" json:"permission"`
	NoteID     uuid.UUID  `gorm:"type:uuid" json:"noteId"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}
