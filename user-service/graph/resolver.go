package graph

import "user-service/internal/user"

type Resolver struct {
	UserService *user.Service
}
