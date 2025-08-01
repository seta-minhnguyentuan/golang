package main

import (
	"user-service/graph"
	"user-service/graph/generated"
	"user-service/internal/user"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=user_service port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&user.User{})

	repo := &user.GormRepository{DB: db}
	service := &user.Service{Repo: repo}

	r := gin.Default()
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{UserService: service}}))

	r.POST("/query", gin.WrapH(srv))
	r.GET("/query", gin.WrapH(playground.Handler("GraphQL", "/query")))

	r.Run(":8080")
}
