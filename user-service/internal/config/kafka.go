package config

import (
	"shared/utils"
	"strings"
	"time"
)

type KafkaConfig struct {
	ServiceName string

	// Kafka
	KafkaBrokers           []string
	KafkaTopicTeams        string
	KafkaTopicTeamActivity string // New topic for team activity events
	KafkaGroupID           string

	// Consumer tuning
	KafkaMinBytes int
	KafkaMaxBytes int
	KafkaMaxWait  time.Duration

	// Producer tuning
	KafkaBatchBytes   int64
	KafkaBatchTimeout time.Duration
}

func LoadKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		ServiceName:            utils.MustEnv("SERVICE_NAME", "user-service"),
		KafkaBrokers:           strings.Split(utils.MustEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaTopicTeams:        utils.MustEnv("KAFKA_TOPIC_TEAMS", "teams"),
		KafkaTopicTeamActivity: utils.MustEnv("KAFKA_TOPIC_TEAM_ACTIVITY", "team.activity"),
		KafkaGroupID:           utils.MustEnv("KAFKA_GROUP_ID", "user-service-consumer"),
		KafkaMinBytes:          utils.AsInt("KAFKA_MIN_BYTES", 10e3), // 10KB
		KafkaMaxBytes:          utils.AsInt("KAFKA_MAX_BYTES", 10e6), // 10MB
		KafkaMaxWait:           utils.AsDuration("KAFKA_MAX_WAIT", 250*time.Millisecond),
		KafkaBatchBytes:        utils.AsInt64("KAFKA_BATCH_BYTES", 1048576), // 1MB
		KafkaBatchTimeout:      utils.AsDuration("KAFKA_BATCH_TIMEOUT", 50*time.Millisecond),
	}
}
