package handlers

import (
	"user-service/internal/services"
)

type Handlers struct {
	TeamHandler *TeamHandler
	UserService services.UserService
}

func NewHandlers(userService services.UserService, teamService services.TeamService) *Handlers {
	return &Handlers{
		TeamHandler: NewTeamHandler(teamService),
		UserService: userService,
	}
}
