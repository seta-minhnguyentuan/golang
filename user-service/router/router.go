package router

import (
	"user-service/graph"
	"user-service/graph/generated"
	"user-service/internal/auth"
	"user-service/internal/team"
	"user-service/internal/user"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

type Router struct {
	UserService *user.Service
	TeamHandler *team.Handler
}

func NewRouter(userService *user.Service, teamService *team.Service) *Router {
	return &Router{
		UserService: userService,
		TeamHandler: team.NewHandler(teamService),
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	userGroup := engine.Group("/user")
	{
		userGroup.POST("/query", r.graphQLHandler())
		userGroup.GET("/query", r.graphQLPlayground())
	}

	teamsGroup := engine.Group("/teams")
	teamsGroup.Use(auth.JWTMiddleware())
	{
		teamsGroup.GET("", r.TeamHandler.GetAllTeams)
		teamsGroup.POST("", r.TeamHandler.CreateTeam)
		teamsGroup.GET("/:teamId", r.TeamHandler.GetTeam)
		teamsGroup.POST("/:teamId/members", r.TeamHandler.AddMember)
		teamsGroup.DELETE("/:teamId/members/:memberId", r.TeamHandler.RemoveMember)
		teamsGroup.POST("/:teamId/managers", r.TeamHandler.AddManager)
		teamsGroup.DELETE("/:teamId/managers/:managerId", r.TeamHandler.RemoveManager)
	}
}

func (r *Router) graphQLPlayground() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/user/query")
	return gin.WrapH(h)
}

func (r *Router) graphQLHandler() gin.HandlerFunc {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{UserService: r.UserService},
	}))
	return gin.WrapH(srv)
}
