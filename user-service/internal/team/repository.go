package team

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, team *Team) (*Team, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Team, error)
	FindAll(ctx context.Context) ([]*Team, error)
	FindMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*TeamMember, error)
	AddMember(ctx context.Context, teamMember *TeamMember) error
	RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error
	FindUserTeams(ctx context.Context, userID uuid.UUID) ([]*Team, error)
	IsUserInTeam(ctx context.Context, teamID, userID uuid.UUID) bool
	IsUserManagerOfTeam(ctx context.Context, teamID, userID uuid.UUID) bool
}

type GormRepository struct {
	DB *gorm.DB
}

func (r *GormRepository) Create(ctx context.Context, team *Team) (*Team, error) {
	if err := r.DB.WithContext(ctx).Create(team).Error; err != nil {
		return nil, err
	}
	return team, nil
}

func (r *GormRepository) FindByID(ctx context.Context, id uuid.UUID) (*Team, error) {
	var team Team
	if err := r.DB.WithContext(ctx).First(&team, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *GormRepository) FindAll(ctx context.Context) ([]*Team, error) {
	var teams []*Team
	if err := r.DB.WithContext(ctx).Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *GormRepository) FindMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*TeamMember, error) {
	var members []*TeamMember
	if err := r.DB.WithContext(ctx).
		Preload("User"). // Preload user information
		Where("team_id = ?", teamID).
		Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *GormRepository) AddMember(ctx context.Context, teamMember *TeamMember) error {
	// Check if user is already in the team
	var count int64
	r.DB.WithContext(ctx).Model(&TeamMember{}).
		Where("team_id = ? AND user_id = ?", teamMember.TeamID, teamMember.UserID).
		Count(&count)

	if count > 0 {
		return errors.New("user is already a member of this team")
	}

	return r.DB.WithContext(ctx).Create(teamMember).Error
}

func (r *GormRepository) RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&TeamMember{}).Error
}

func (r *GormRepository) FindUserTeams(ctx context.Context, userID uuid.UUID) ([]*Team, error) {
	var teams []*Team
	if err := r.DB.WithContext(ctx).
		Table("teams").
		Joins("JOIN team_members ON teams.id = team_members.team_id").
		Where("team_members.user_id = ?", userID).
		Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *GormRepository) IsUserInTeam(ctx context.Context, teamID, userID uuid.UUID) bool {
	var count int64
	r.DB.WithContext(ctx).Model(&TeamMember{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Count(&count)
	return count > 0
}

func (r *GormRepository) IsUserManagerOfTeam(ctx context.Context, teamID, userID uuid.UUID) bool {
	var count int64
	r.DB.WithContext(ctx).Model(&TeamMember{}).
		Where("team_id = ? AND user_id = ? AND role = ?", teamID, userID, "manager").
		Count(&count)
	return count > 0
}
