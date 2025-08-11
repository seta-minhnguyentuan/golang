package models

import (
	"time"

	"github.com/google/uuid"
)

type Permission string

const (
	PermissionRead  Permission = "read"
	PermissionWrite Permission = "write"
	PermissionOwner Permission = "owner"
)

type Sharing struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null" json:"userId"`
	Permission Permission `gorm:"type:varchar(16);not null;check:permission IN ('read','write','owner')" json:"permission"`
	FolderID   *uuid.UUID `gorm:"type:uuid" json:"folderId,omitempty"`
	NoteID     *uuid.UUID `gorm:"type:uuid" json:"noteId,omitempty"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}
