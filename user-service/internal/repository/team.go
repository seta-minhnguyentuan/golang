package repository

import (
	"context"
	"errors"
	"user-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository interface {
	Create(ctx context.Context, team *models.Team) (*models.Team, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Team, error)
	FindAll(ctx context.Context) ([]*models.Team, error)
	FindMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*models.TeamMember, error)
	AddMember(ctx context.Context, teamMember *models.TeamMember) error
	RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error
	FindUserTeams(ctx context.Context, userID uuid.UUID) ([]*models.Team, error)
	IsUserInTeam(ctx context.Context, teamID, userID uuid.UUID) bool
	IsUserManagerOfTeam(ctx context.Context, teamID, userID uuid.UUID) bool
}

type GormTeamRepository struct {
	DB *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &GormTeamRepository{DB: db}
}

func (r *GormTeamRepository) Create(ctx context.Context, team *models.Team) (*models.Team, error) {
	if err := r.DB.WithContext(ctx).Create(team).Error; err != nil {
		return nil, err
	}
	return team, nil
}

func (r *GormTeamRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Team, error) {
	var team models.Team
	if err := r.DB.WithContext(ctx).First(&team, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *GormTeamRepository) FindAll(ctx context.Context) ([]*models.Team, error) {
	var teams []*models.Team
	if err := r.DB.WithContext(ctx).Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *GormTeamRepository) FindMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*models.TeamMember, error) {
	var members []*models.TeamMember
	if err := r.DB.WithContext(ctx).
		Preload("User"). // Preload user information
		Where("team_id = ?", teamID).
		Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *GormTeamRepository) AddMember(ctx context.Context, teamMember *models.TeamMember) error {
	// Check if user is already in the team
	var count int64
	r.DB.WithContext(ctx).Model(&models.TeamMember{}).
		Where("team_id = ? AND user_id = ?", teamMember.TeamID, teamMember.UserID).
		Count(&count)

	if count > 0 {
		return errors.New("user is already a member of this team")
	}

	return r.DB.WithContext(ctx).Create(teamMember).Error
}

func (r *GormTeamRepository) RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&models.TeamMember{}).Error
}

func (r *GormTeamRepository) FindUserTeams(ctx context.Context, userID uuid.UUID) ([]*models.Team, error) {
	var teams []*models.Team
	if err := r.DB.WithContext(ctx).
		Table("teams").
		Joins("JOIN team_members ON teams.id = team_members.team_id").
		Where("team_members.user_id = ?", userID).
		Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *GormTeamRepository) IsUserInTeam(ctx context.Context, teamID, userID uuid.UUID) bool {
	var count int64
	r.DB.WithContext(ctx).Model(&models.TeamMember{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Count(&count)
	return count > 0
}

func (r *GormTeamRepository) IsUserManagerOfTeam(ctx context.Context, teamID, userID uuid.UUID) bool {
	var count int64
	r.DB.WithContext(ctx).Model(&models.TeamMember{}).
		Where("team_id = ? AND user_id = ? AND role = ?", teamID, userID, "manager").
		Count(&count)
	return count > 0
}
