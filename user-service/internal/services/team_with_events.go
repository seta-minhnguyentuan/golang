package services

import (
	"context"
	"time"
	"user-service/internal/kafka"
	"user-service/internal/models"

	"shared/pkg/log"

	"github.com/google/uuid"
)

// TeamServiceWithEvents wraps the existing TeamService and adds Kafka event publishing
type TeamServiceWithEvents struct {
	baseService TeamService
	producer    kafka.Producer
	topicName   string
}

func NewTeamServiceWithEvents(baseService TeamService, producer kafka.Producer, topicName string) TeamService {
	return &TeamServiceWithEvents{
		baseService: baseService,
		producer:    producer,
		topicName:   topicName,
	}
}

func (s *TeamServiceWithEvents) CreateTeam(ctx context.Context, req models.CreateTeamRequest, creatorID uuid.UUID) (*models.TeamResponse, error) {
	// Call the base service to create the team
	teamResponse, err := s.baseService.CreateTeam(ctx, req, creatorID)
	if err != nil {
		return nil, err
	}

	// Publish TEAM_CREATED event
	event := kafka.TeamActivityEvent{
		EventType:   kafka.EventTypeTeamCreated,
		TeamID:      teamResponse.ID.String(),
		PerformedBy: creatorID.String(),
		TeamName:    &teamResponse.TeamName,
		Timestamp:   time.Now(),
	}

	if err := s.publishEvent(ctx, event); err != nil {
		log.Error.Printf("Failed to publish TEAM_CREATED event: %v", err)
		// Don't fail the operation if event publishing fails
	}

	// Publish MANAGER_ADDED events for all managers (including creator)
	for _, manager := range teamResponse.Managers {
		managerEvent := kafka.TeamActivityEvent{
			EventType:    kafka.EventTypeManagerAdded,
			TeamID:       teamResponse.ID.String(),
			PerformedBy:  creatorID.String(),
			TargetUserID: stringPtr(manager.UserID.String()),
			Timestamp:    time.Now(),
		}
		if err := s.publishEvent(ctx, managerEvent); err != nil {
			log.Error.Printf("Failed to publish MANAGER_ADDED event: %v", err)
		}
	}

	// Publish MEMBER_ADDED events for all members
	for _, member := range teamResponse.Members {
		memberEvent := kafka.TeamActivityEvent{
			EventType:    kafka.EventTypeMemberAdded,
			TeamID:       teamResponse.ID.String(),
			PerformedBy:  creatorID.String(),
			TargetUserID: stringPtr(member.UserID.String()),
			Timestamp:    time.Now(),
		}
		if err := s.publishEvent(ctx, memberEvent); err != nil {
			log.Error.Printf("Failed to publish MEMBER_ADDED event: %v", err)
		}
	}

	return teamResponse, nil
}

func (s *TeamServiceWithEvents) GetTeamByID(ctx context.Context, teamID uuid.UUID) (*models.TeamResponse, error) {
	return s.baseService.GetTeamByID(ctx, teamID)
}

func (s *TeamServiceWithEvents) AddMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, requestorID uuid.UUID) error {
	// Call the base service to add the member
	err := s.baseService.AddMember(ctx, teamID, userID, requestorID)
	if err != nil {
		return err
	}

	// Publish MEMBER_ADDED event
	event := kafka.TeamActivityEvent{
		EventType:    kafka.EventTypeMemberAdded,
		TeamID:       teamID.String(),
		PerformedBy:  requestorID.String(),
		TargetUserID: stringPtr(userID.String()),
		Timestamp:    time.Now(),
	}

	if err := s.publishEvent(ctx, event); err != nil {
		log.Error.Printf("Failed to publish MEMBER_ADDED event: %v", err)
		// Don't fail the operation if event publishing fails
	}

	return nil
}

func (s *TeamServiceWithEvents) RemoveMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, requestorID uuid.UUID) error {
	// Call the base service to remove the member
	err := s.baseService.RemoveMember(ctx, teamID, userID, requestorID)
	if err != nil {
		return err
	}

	// Publish MEMBER_REMOVED event
	event := kafka.TeamActivityEvent{
		EventType:    kafka.EventTypeMemberRemoved,
		TeamID:       teamID.String(),
		PerformedBy:  requestorID.String(),
		TargetUserID: stringPtr(userID.String()),
		Timestamp:    time.Now(),
	}

	if err := s.publishEvent(ctx, event); err != nil {
		log.Error.Printf("Failed to publish MEMBER_REMOVED event: %v", err)
		// Don't fail the operation if event publishing fails
	}

	return nil
}

func (s *TeamServiceWithEvents) AddManager(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, requestorID uuid.UUID) error {
	// Call the base service to add the manager
	err := s.baseService.AddManager(ctx, teamID, userID, requestorID)
	if err != nil {
		return err
	}

	// Publish MANAGER_ADDED event
	event := kafka.TeamActivityEvent{
		EventType:    kafka.EventTypeManagerAdded,
		TeamID:       teamID.String(),
		PerformedBy:  requestorID.String(),
		TargetUserID: stringPtr(userID.String()),
		Timestamp:    time.Now(),
	}

	if err := s.publishEvent(ctx, event); err != nil {
		log.Error.Printf("Failed to publish MANAGER_ADDED event: %v", err)
		// Don't fail the operation if event publishing fails
	}

	return nil
}

func (s *TeamServiceWithEvents) RemoveManager(ctx context.Context, teamID uuid.UUID, managerID uuid.UUID, requestorID uuid.UUID) error {
	// Call the base service to remove the manager
	err := s.baseService.RemoveManager(ctx, teamID, managerID, requestorID)
	if err != nil {
		return err
	}

	// Publish MANAGER_REMOVED event
	event := kafka.TeamActivityEvent{
		EventType:    kafka.EventTypeManagerRemoved,
		TeamID:       teamID.String(),
		PerformedBy:  requestorID.String(),
		TargetUserID: stringPtr(managerID.String()),
		Timestamp:    time.Now(),
	}

	if err := s.publishEvent(ctx, event); err != nil {
		log.Error.Printf("Failed to publish MANAGER_REMOVED event: %v", err)
		// Don't fail the operation if event publishing fails
	}

	return nil
}

func (s *TeamServiceWithEvents) GetAllTeams(ctx context.Context, requestorID uuid.UUID) ([]*models.TeamResponse, error) {
	return s.baseService.GetAllTeams(ctx, requestorID)
}

// publishEvent publishes a team activity event to Kafka
func (s *TeamServiceWithEvents) publishEvent(ctx context.Context, event kafka.TeamActivityEvent) error {
	key := []byte(event.TeamID) // Use team ID as the key for partitioning
	return s.producer.Publish(ctx, s.topicName, key, event)
}

// stringPtr is a helper function to create a pointer to a string
func stringPtr(s string) *string {
	return &s
}
