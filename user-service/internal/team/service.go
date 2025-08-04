package team

import (
	"context"
	"errors"
	"time"
	"user-service/internal/user"

	"github.com/google/uuid"
)

type Service struct {
	TeamRepo Repository
	UserRepo user.Repository
}

func (s *Service) CreateTeam(ctx context.Context, req CreateTeamRequest, creatorID uuid.UUID) (*TeamResponse, error) {
	// Verify that the creator is a manager
	creator, err := s.UserRepo.FindByID(ctx, creatorID)
	if err != nil {
		return nil, errors.New("creator not found")
	}

	if creator.Role != "manager" {
		return nil, errors.New("only managers can create teams")
	}

	// Create the team
	team := &Team{
		ID:        uuid.New(),
		TeamName:  req.TeamName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdTeam, err := s.TeamRepo.Create(ctx, team)
	if err != nil {
		return nil, err
	}

	// Add the creator as a manager
	creatorMember := &TeamMember{
		ID:       uuid.New(),
		TeamID:   createdTeam.ID,
		UserID:   creatorID,
		Role:     "manager",
		JoinedAt: time.Now(),
	}

	if err := s.TeamRepo.AddMember(ctx, creatorMember); err != nil {
		return nil, err
	}

	// Add other managers
	for _, manager := range req.Managers {
		managerUUID, err := uuid.Parse(manager.UserID)
		if err != nil {
			continue // Skip invalid UUIDs
		}

		// Skip if it's the same as creator
		if managerUUID == creatorID {
			continue
		}

		// Verify user exists and is a manager
		managerUser, err := s.UserRepo.FindByID(ctx, managerUUID)
		if err != nil || managerUser.Role != "manager" {
			continue // Skip if user doesn't exist or isn't a manager
		}

		teamMember := &TeamMember{
			ID:       uuid.New(),
			TeamID:   createdTeam.ID,
			UserID:   managerUUID,
			Role:     "manager",
			JoinedAt: time.Now(),
		}
		s.TeamRepo.AddMember(ctx, teamMember)
	}

	// Add members
	for _, member := range req.Members {
		memberUUID, err := uuid.Parse(member.UserID)
		if err != nil {
			continue // Skip invalid UUIDs
		}

		// Skip if it's the same as creator
		if memberUUID == creatorID {
			continue
		}

		// Verify user exists
		_, err = s.UserRepo.FindByID(ctx, memberUUID)
		if err != nil {
			continue // Skip if user doesn't exist
		}

		teamMember := &TeamMember{
			ID:       uuid.New(),
			TeamID:   createdTeam.ID,
			UserID:   memberUUID,
			Role:     "member",
			JoinedAt: time.Now(),
		}
		s.TeamRepo.AddMember(ctx, teamMember)
	}

	return s.GetTeamByID(ctx, createdTeam.ID)
}

func (s *Service) GetTeamByID(ctx context.Context, teamID uuid.UUID) (*TeamResponse, error) {
	team, err := s.TeamRepo.FindByID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	members, err := s.TeamRepo.FindMembersByTeamID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	var managers []TeamMemberResponse
	var teamMembers []TeamMemberResponse

	for _, member := range members {
		// If User is preloaded, use it; otherwise fetch separately
		var memberUser *user.User
		if member.User.ID != uuid.Nil {
			memberUser = &member.User
		} else {
			memberUser, err = s.UserRepo.FindByID(ctx, member.UserID)
			if err != nil {
				continue // Skip if user not found
			}
		}

		memberResponse := TeamMemberResponse{
			UserID:   memberUser.ID,
			UserName: memberUser.Username,
			Email:    memberUser.Email,
			Role:     member.Role,
			JoinedAt: member.JoinedAt,
		}

		if member.Role == "manager" {
			managers = append(managers, memberResponse)
		} else {
			teamMembers = append(teamMembers, memberResponse)
		}
	}

	return &TeamResponse{
		ID:        team.ID,
		TeamName:  team.TeamName,
		Managers:  managers,
		Members:   teamMembers,
		CreatedAt: team.CreatedAt,
		UpdatedAt: team.UpdatedAt,
	}, nil
}

func (s *Service) AddMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, requestorID uuid.UUID) error {
	// Check if requestor is a manager of the team
	if !s.TeamRepo.IsUserManagerOfTeam(ctx, teamID, requestorID) {
		return errors.New("only team managers can add members")
	}

	// Check if user exists
	_, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if user is already in the team
	if s.TeamRepo.IsUserInTeam(ctx, teamID, userID) {
		return errors.New("user is already a member of this team")
	}

	teamMember := &TeamMember{
		ID:       uuid.New(),
		TeamID:   teamID,
		UserID:   userID,
		Role:     "member",
		JoinedAt: time.Now(),
	}

	return s.TeamRepo.AddMember(ctx, teamMember)
}

func (s *Service) RemoveMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, requestorID uuid.UUID) error {
	// Check if requestor is a manager of the team
	if !s.TeamRepo.IsUserManagerOfTeam(ctx, teamID, requestorID) {
		return errors.New("only team managers can remove members")
	}

	// Check if user is in the team
	if !s.TeamRepo.IsUserInTeam(ctx, teamID, userID) {
		return errors.New("user is not a member of this team")
	}

	return s.TeamRepo.RemoveMember(ctx, teamID, userID)
}

func (s *Service) AddManager(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, requestorID uuid.UUID) error {
	// Check if requestor is a manager of the team
	if !s.TeamRepo.IsUserManagerOfTeam(ctx, teamID, requestorID) {
		return errors.New("only team managers can add other managers")
	}

	// Check if user exists and is a manager role
	user, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Role != "manager" {
		return errors.New("user must have manager role to be added as team manager")
	}

	// Check if user is already in the team
	if s.TeamRepo.IsUserInTeam(ctx, teamID, userID) {
		return errors.New("user is already a member of this team")
	}

	teamMember := &TeamMember{
		ID:       uuid.New(),
		TeamID:   teamID,
		UserID:   userID,
		Role:     "manager",
		JoinedAt: time.Now(),
	}

	return s.TeamRepo.AddMember(ctx, teamMember)
}

func (s *Service) RemoveManager(ctx context.Context, teamID uuid.UUID, managerID uuid.UUID, requestorID uuid.UUID) error {
	// Check if requestor is a manager of the team
	if !s.TeamRepo.IsUserManagerOfTeam(ctx, teamID, requestorID) {
		return errors.New("only team managers can remove other managers")
	}

	// Check if the manager to be removed is in the team as a manager
	if !s.TeamRepo.IsUserManagerOfTeam(ctx, teamID, managerID) {
		return errors.New("user is not a manager of this team")
	}

	return s.TeamRepo.RemoveMember(ctx, teamID, managerID)
}

func (s *Service) GetAllTeams(ctx context.Context, requestorID uuid.UUID) ([]*TeamResponse, error) {
	requestor, err := s.UserRepo.FindByID(ctx, requestorID)
	if err != nil {
		return nil, errors.New("requestor not found")
	}

	var teams []*Team

	if requestor.Role == "manager" {
		teams, err = s.TeamRepo.FindAll(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		teams, err = s.TeamRepo.FindUserTeams(ctx, requestorID)
		if err != nil {
			return nil, err
		}
	}

	var teamResponses []*TeamResponse
	for _, team := range teams {
		teamResponse, err := s.GetTeamByID(ctx, team.ID)
		if err != nil {
			continue
		}
		teamResponses = append(teamResponses, teamResponse)
	}

	return teamResponses, nil
}
