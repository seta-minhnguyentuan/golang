package kafka

import "time"

type Key []byte
type Value []byte

// Team Activity Event Types
const (
	EventTypeTeamCreated    = "TEAM_CREATED"
	EventTypeMemberAdded    = "MEMBER_ADDED"
	EventTypeMemberRemoved  = "MEMBER_REMOVED"
	EventTypeManagerAdded   = "MANAGER_ADDED"
	EventTypeManagerRemoved = "MANAGER_REMOVED"
)

// Team Activity Event as specified in kafka_redis.md
type TeamActivityEvent struct {
	EventType    string    `json:"eventType"`
	TeamID       string    `json:"teamId"`
	PerformedBy  string    `json:"performedBy"`
	TargetUserID *string   `json:"targetUserId,omitempty"` // nil for TEAM_CREATED
	TeamName     *string   `json:"teamName,omitempty"`     // only for TEAM_CREATED
	Timestamp    time.Time `json:"timestamp"`
}

// Legacy type for backward compatibility
type TeamCreated struct {
	TeamID   string `json:"team_id"`
	TeamName string `json:"team_name"`
}
