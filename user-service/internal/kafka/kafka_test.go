package kafka

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTeamActivityEvent_JSON(t *testing.T) {
	// Test event serialization/deserialization
	teamID := uuid.New().String()
	userID := uuid.New().String()
	targetUserID := uuid.New().String()

	event := TeamActivityEvent{
		EventType:    EventTypeMemberAdded,
		TeamID:       teamID,
		PerformedBy:  userID,
		TargetUserID: &targetUserID,
		Timestamp:    time.Now(),
	}

	// Marshal to JSON
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}

	// Unmarshal from JSON
	var parsed TeamActivityEvent
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Verify fields
	if parsed.EventType != event.EventType {
		t.Errorf("Expected eventType %s, got %s", event.EventType, parsed.EventType)
	}
	if parsed.TeamID != event.TeamID {
		t.Errorf("Expected teamId %s, got %s", event.TeamID, parsed.TeamID)
	}
	if parsed.PerformedBy != event.PerformedBy {
		t.Errorf("Expected performedBy %s, got %s", event.PerformedBy, parsed.PerformedBy)
	}
	if parsed.TargetUserID == nil || *parsed.TargetUserID != *event.TargetUserID {
		t.Errorf("Expected targetUserId %s, got %v", *event.TargetUserID, parsed.TargetUserID)
	}
}

func TestTeamActivityEventHandler_HandleEvent(t *testing.T) {
	handler := NewTeamActivityEventHandler()
	ctx := context.Background()

	// Test TEAM_CREATED event
	teamCreatedEvent := TeamActivityEvent{
		EventType:   EventTypeTeamCreated,
		TeamID:      uuid.New().String(),
		PerformedBy: uuid.New().String(),
		TeamName:    stringPtr("Test Team"),
		Timestamp:   time.Now(),
	}

	eventData, _ := json.Marshal(teamCreatedEvent)
	key := []byte(teamCreatedEvent.TeamID)

	err := handler.HandleEvent(ctx, key, eventData)
	if err != nil {
		t.Errorf("Expected no error for TEAM_CREATED event, got: %v", err)
	}

	// Test MEMBER_ADDED event
	memberAddedEvent := TeamActivityEvent{
		EventType:    EventTypeMemberAdded,
		TeamID:       uuid.New().String(),
		PerformedBy:  uuid.New().String(),
		TargetUserID: stringPtr(uuid.New().String()),
		Timestamp:    time.Now(),
	}

	eventData, _ = json.Marshal(memberAddedEvent)
	key = []byte(memberAddedEvent.TeamID)

	err = handler.HandleEvent(ctx, key, eventData)
	if err != nil {
		t.Errorf("Expected no error for MEMBER_ADDED event, got: %v", err)
	}

	// Test invalid event (missing targetUserId for MEMBER_ADDED)
	invalidEvent := TeamActivityEvent{
		EventType:   EventTypeMemberAdded,
		TeamID:      uuid.New().String(),
		PerformedBy: uuid.New().String(),
		// TargetUserID is nil
		Timestamp: time.Now(),
	}

	eventData, _ = json.Marshal(invalidEvent)
	key = []byte(invalidEvent.TeamID)

	err = handler.HandleEvent(ctx, key, eventData)
	if err == nil {
		t.Error("Expected error for MEMBER_ADDED event without targetUserId")
	}
}

func TestEventConstants(t *testing.T) {
	// Verify event type constants match specification
	expectedEvents := []string{
		"TEAM_CREATED",
		"MEMBER_ADDED",
		"MEMBER_REMOVED",
		"MANAGER_ADDED",
		"MANAGER_REMOVED",
	}

	actualEvents := []string{
		EventTypeTeamCreated,
		EventTypeMemberAdded,
		EventTypeMemberRemoved,
		EventTypeManagerAdded,
		EventTypeManagerRemoved,
	}

	for i, expected := range expectedEvents {
		if actualEvents[i] != expected {
			t.Errorf("Expected event type %s, got %s", expected, actualEvents[i])
		}
	}
}

// Helper function for tests
func stringPtr(s string) *string {
	return &s
}
