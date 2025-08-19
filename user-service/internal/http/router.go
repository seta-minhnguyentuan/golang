package httpserver

import (
	"net/http"
	"shared/middlewares"
	"user-service/graph"
	"user-service/graph/generated"
	"user-service/internal/handlers"
	"user-service/internal/services"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	UserService services.UserService
	TeamService services.TeamService
}

func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:5174", "http://localhost:4173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // 12 hours
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	h := handlers.NewHandlers(deps.UserService, deps.TeamService)

	// GraphQL routes for user service
	userGroup := r.Group("/user")
	// userGroup.Use(middlewares.AuthMiddleware())
	{
		userGroup.POST("/query", graphQLHandler(deps.UserService))
		userGroup.GET("/query", graphQLPlayground())
	}

	// REST API routes for teams
	teamsGroup := r.Group("/teams")
	teamsGroup.Use(middlewares.AuthMiddleware())
	{
		teamsGroup.GET("", h.TeamHandler.GetAllTeams)
		teamsGroup.POST("", h.TeamHandler.CreateTeam)
		teamsGroup.GET("/:teamId", h.TeamHandler.GetTeam)
		teamsGroup.POST("/:teamId/members", h.TeamHandler.AddMember)
		teamsGroup.DELETE("/:teamId/members/:memberId", h.TeamHandler.RemoveMember)
		teamsGroup.POST("/:teamId/managers", h.TeamHandler.AddManager)
		teamsGroup.DELETE("/:teamId/managers/:managerId", h.TeamHandler.RemoveManager)
	}

	return r
}

func graphQLPlayground() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/user/query")
	return gin.WrapH(h)
}

func graphQLHandler(userService services.UserService) gin.HandlerFunc {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{UserService: userService},
	}))
	return gin.WrapH(srv)
}
