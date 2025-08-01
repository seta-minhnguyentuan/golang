package graph

import (
	"context"
	"fmt"

	"user-service/graph/model"
	"user-service/internal/auth"
	"user-service/internal/user"
)

// Helper function to convert internal user model to GraphQL model
func toGraphQLUser(user *user.User) *model.User {
	return &model.User{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
}

// CreateUser creates a new user
func (r *mutationResolver) CreateUser(ctx context.Context, username string, email string, password string, role string) (*model.User, error) {
	// Validate input parameters
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	// Validate role
	if role != "manager" && role != "member" {
		return nil, fmt.Errorf("invalid role: must be either 'manager' or 'member'")
	}

	// Check if user with email already exists
	existingUser, err := r.UserService.Repo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email '%s' already exists", email)
	}

	// Use the service layer to create user (includes password hashing)
	user, err := r.UserService.CreateUser(ctx, username, email, password, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Convert internal user model to GraphQL model
	return toGraphQLUser(user), nil
}

// Login authenticates a user and returns a JWT token
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.AuthPayload, error) {
	// Use the service layer to authenticate user
	user, err := r.UserService.Login(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID.String(), user.Email, user.Role, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.AuthPayload{
		Token: token,
		User:  toGraphQLUser(user),
	}, nil
}

// Logout handles user logout (for now just returns true as JWT is stateless)
func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	// In a stateless JWT system, logout is typically handled client-side
	// by removing the token. For more advanced scenarios, you might want
	// to implement a token blacklist.
	return true, nil
}

// FetchUsers returns all users
func (r *queryResolver) FetchUsers(ctx context.Context) ([]*model.User, error) {
	users, err := r.UserService.FetchUsers(ctx)
	if err != nil {
		return nil, err
	}

	var graphqlUsers []*model.User
	for _, user := range users {
		graphqlUsers = append(graphqlUsers, toGraphQLUser(user))
	}

	return graphqlUsers, nil
}
