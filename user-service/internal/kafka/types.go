package kafka

type Key []byte
type Value []byte

type TeamCreated struct {
	TeamID   string `json:"team_id"`
	TeamName string `json:"team_name"`
}
