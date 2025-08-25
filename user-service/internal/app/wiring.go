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
	producer := kafka.NewProducer(cfg)
	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	teamService := services.NewTeamService(teamRepo, userRepo)

	consumer := kafka.NewConsumer(cfg, cfg.KafkaTopicTeams, teamService.CreateTeam())

	return &Components{
		Cfg:      cfg,
		Producer: producer,
		Teams:    teamService,
		Consumer: consumer,
	}
}

func (c *Components) PublishTestTeam(ctx context.Context, id string) error {
	return c.Producer.Publish(ctx, c.Cfg.KafkaTopicTeams, []byte(id), kafka.TeamCreated{
		TeamID:   id,
		TeamName: "Test Team",
	})
}
