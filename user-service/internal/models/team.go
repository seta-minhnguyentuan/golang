package models

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TeamName  string    `gorm:"not null" json:"teamName"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	TeamMembers []TeamMember `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE" json:"teamMembers,omitempty"`
}

type TeamMember struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TeamID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_team_user" json:"teamId"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_team_user" json:"userId"`
	Role     string    `gorm:"type:VARCHAR(10);not null" json:"role"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joinedAt"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (TeamMember) TableName() string {
	return "team_members"
}

type CreateTeamRequest struct {
	TeamName string              `json:"teamName" binding:"required"`
	Managers []TeamMemberRequest `json:"managers"`
	Members  []TeamMemberRequest `json:"members"`
}

type TeamMemberRequest struct {
	UserID   string `json:"userId" binding:"required"`
	UserName string `json:"userName"`
}

type AddMemberRequest struct {
	UserID string `json:"userId" binding:"required"`
}

type TeamResponse struct {
	ID        uuid.UUID            `json:"id"`
	TeamName  string               `json:"teamName"`
	Managers  []TeamMemberResponse `json:"managers"`
	Members   []TeamMemberResponse `json:"members"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
}

type TeamMemberResponse struct {
	UserID   uuid.UUID `json:"userId"`
	UserName string    `json:"userName"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joinedAt"`
}
