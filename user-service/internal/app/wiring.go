package app

import (
	"context"
	"user-service/internal/config"
	"user-service/internal/database"
	"user-service/internal/kafka"
	"user-service/internal/repository"
	"user-service/internal/services"
)

type Components struct {
	Cfg      *config.KafkaConfig
	Producer kafka.Producer
	Consumer *kafka.Consumer
	Users    services.UserService
	Teams    services.TeamService
}

func Wire(cfg *config.KafkaConfig) *Components {
	dbCfg := config.LoadDB()

	db, err := database.Connect(*dbCfg)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Initialize repositories
	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize Kafka producer
	producer := kafka.NewProducer(cfg)

	// Initialize base team service
	baseTeamService := services.NewTeamService(teamRepo, userRepo)

	// Wrap base team service with event publishing
	teamService := services.NewTeamServiceWithEvents(baseTeamService, producer, cfg.KafkaTopicTeamActivity)

	// Initialize user service
	userService := services.NewUserService(userRepo)

	// Initialize event handler and consumer
	eventHandler := kafka.NewTeamActivityEventHandler()
	consumer := kafka.NewConsumer(cfg, cfg.KafkaTopicTeamActivity, eventHandler.HandleEvent)

	return &Components{
		Cfg:      cfg,
		Producer: producer,
		Teams:    teamService,
		Users:    userService,
		Consumer: consumer,
	}
}

func (c *Components) PublishTestTeam(ctx context.Context, id string) error {
	return c.Producer.Publish(ctx, c.Cfg.KafkaTopicTeams, []byte(id), kafka.TeamCreated{
		TeamID:   id,
		TeamName: "Test Team",
	})
}
