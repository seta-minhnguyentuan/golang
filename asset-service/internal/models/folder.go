package models

import (
	"time"

	"github.com/google/uuid"
)

type Folder struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string          `gorm:"not null" json:"folderName"`
	Notes     []Note          `gorm:"foreignKey:FolderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"notes"`
	Sharings  []FolderSharing `gorm:"foreignKey:FolderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"sharings"`
	OwnerID   uuid.UUID       `gorm:"type:uuid;not null" json:"ownerId"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`
	CreatedBy uuid.UUID       `gorm:"type:uuid;not null" json:"createdBy"`
	UpdatedBy uuid.UUID       `gorm:"type:uuid;not null" json:"updatedBy"`
}
