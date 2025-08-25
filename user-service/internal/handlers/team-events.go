package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"shared/pkg/log"
	"user-service/internal/kafka"
)

// TeamActivityEventHandler handles team activity events
type TeamActivityEventHandler struct {
	// Add dependencies here for logging, database, etc.
	// For example:
	// auditRepo repository.AuditRepository
	// cacheService cache.CacheService
}

func NewTeamActivityEventHandler() *TeamActivityEventHandler {
	return &TeamActivityEventHandler{}
}

// HandleEvent processes a team activity event
func (h *TeamActivityEventHandler) HandleEvent(ctx context.Context, key []byte, value []byte) error {
	var event kafka.TeamActivityEvent
	if err := json.Unmarshal(value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal team activity event: %w", err)
	}

	log.Info.Printf("Processing team activity event: %s for team %s", event.EventType, event.TeamID)

	switch event.EventType {
	case kafka.EventTypeTeamCreated:
		return h.handleTeamCreated(ctx, event)
	case kafka.EventTypeMemberAdded:
		return h.handleMemberAdded(ctx, event)
	case kafka.EventTypeMemberRemoved:
		return h.handleMemberRemoved(ctx, event)
	case kafka.EventTypeManagerAdded:
		return h.handleManagerAdded(ctx, event)
	case kafka.EventTypeManagerRemoved:
		return h.handleManagerRemoved(ctx, event)
	default:
		log.Info.Printf("Unknown event type: %s", event.EventType)
		return nil // Don't fail on unknown events
	}
}

func (h *TeamActivityEventHandler) handleTeamCreated(ctx context.Context, event kafka.TeamActivityEvent) error {
	log.Info.Printf("Team created: %s by user %s", event.TeamID, event.PerformedBy)

	// Example implementations:
	// 1. Store audit log
	// h.auditRepo.CreateAuditLog(ctx, &models.AuditLog{
	//     EventType: "TEAM_CREATED",
	//     TeamID: event.TeamID,
	//     UserID: event.PerformedBy,
	//     Timestamp: event.Timestamp,
	//     Details: map[string]interface{}{
	//         "teamName": *event.TeamName,
	//     },
	// })

	// 2. Update cache
	// h.cacheService.InvalidateTeamCache(ctx, event.TeamID)

	// 3. Send notifications
	// h.notificationService.NotifyTeamCreated(ctx, event.TeamID, event.PerformedBy)

	return nil
}

func (h *TeamActivityEventHandler) handleMemberAdded(ctx context.Context, event kafka.TeamActivityEvent) error {
	if event.TargetUserID == nil {
		return fmt.Errorf("targetUserId is required for MEMBER_ADDED event")
	}

	log.Info.Printf("Member %s added to team %s by user %s", *event.TargetUserID, event.TeamID, event.PerformedBy)

	// Example implementations:
	// 1. Update team member cache
	// h.cacheService.AddTeamMember(ctx, event.TeamID, *event.TargetUserID)

	// 2. Store audit log
	// h.auditRepo.CreateAuditLog(ctx, &models.AuditLog{
	//     EventType: "MEMBER_ADDED",
	//     TeamID: event.TeamID,
	//     UserID: event.PerformedBy,
	//     TargetUserID: event.TargetUserID,
	//     Timestamp: event.Timestamp,
	// })

	// 3. Send welcome notification to new member
	// h.notificationService.NotifyMemberAdded(ctx, event.TeamID, *event.TargetUserID)

	return nil
}

func (h *TeamActivityEventHandler) handleMemberRemoved(ctx context.Context, event kafka.TeamActivityEvent) error {
	if event.TargetUserID == nil {
		return fmt.Errorf("targetUserId is required for MEMBER_REMOVED event")
	}

	log.Info.Printf("Member %s removed from team %s by user %s", *event.TargetUserID, event.TeamID, event.PerformedBy)

	// Example implementations:
	// 1. Update team member cache
	// h.cacheService.RemoveTeamMember(ctx, event.TeamID, *event.TargetUserID)

	// 2. Store audit log
	// h.auditRepo.CreateAuditLog(ctx, &models.AuditLog{
	//     EventType: "MEMBER_REMOVED",
	//     TeamID: event.TeamID,
	//     UserID: event.PerformedBy,
	//     TargetUserID: event.TargetUserID,
	//     Timestamp: event.Timestamp,
	// })

	// 3. Clean up user's access to team resources
	// h.accessControlService.RevokeTeamAccess(ctx, event.TeamID, *event.TargetUserID)

	return nil
}

func (h *TeamActivityEventHandler) handleManagerAdded(ctx context.Context, event kafka.TeamActivityEvent) error {
	if event.TargetUserID == nil {
		return fmt.Errorf("targetUserId is required for MANAGER_ADDED event")
	}

	log.Info.Printf("Manager %s added to team %s by user %s", *event.TargetUserID, event.TeamID, event.PerformedBy)

	// Example implementations:
	// 1. Update team manager cache
	// h.cacheService.AddTeamManager(ctx, event.TeamID, *event.TargetUserID)

	// 2. Store audit log
	// h.auditRepo.CreateAuditLog(ctx, &models.AuditLog{
	//     EventType: "MANAGER_ADDED",
	//     TeamID: event.TeamID,
	//     UserID: event.PerformedBy,
	//     TargetUserID: event.TargetUserID,
	//     Timestamp: event.Timestamp,
	// })

	// 3. Grant manager permissions
	// h.accessControlService.GrantManagerAccess(ctx, event.TeamID, *event.TargetUserID)

	return nil
}

func (h *TeamActivityEventHandler) handleManagerRemoved(ctx context.Context, event kafka.TeamActivityEvent) error {
	if event.TargetUserID == nil {
		return fmt.Errorf("targetUserId is required for MANAGER_REMOVED event")
	}

	log.Info.Printf("Manager %s removed from team %s by user %s", *event.TargetUserID, event.TeamID, event.PerformedBy)

	// Example implementations:
	// 1. Update team manager cache
	// h.cacheService.RemoveTeamManager(ctx, event.TeamID, *event.TargetUserID)

	// 2. Store audit log
	// h.auditRepo.CreateAuditLog(ctx, &models.AuditLog{
	//     EventType: "MANAGER_REMOVED",
	//     TeamID: event.TeamID,
	//     UserID: event.PerformedBy,
	//     TargetUserID: event.TargetUserID,
	//     Timestamp: event.Timestamp,
	// })

	// 3. Revoke manager permissions (but keep member access if still a member)
	// h.accessControlService.RevokeManagerAccess(ctx, event.TeamID, *event.TargetUserID)

	return nil
}
