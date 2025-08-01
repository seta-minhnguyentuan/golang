package graph

import (
	"user-service/graph/generated"
	"user-service/internal/user"
)

type Resolver struct {
	UserService *user.Service
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }