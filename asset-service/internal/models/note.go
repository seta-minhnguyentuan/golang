package models

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"type:string" json:"noteName"`
	Content   string    `gorm:"type:string" json:"noteContent"`
	FolderID  uuid.UUID `gorm:"type:uuid;not null" json:"folderId"`
	Sharings  []Sharing `gorm:"foreignKey:NoteID" json:"sharings"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
