# User Service - GraphQL API

## Overview
This is a GraphQL-based user management service built with Go, GORM, and gqlgen. It provides user authentication, authorization, and user management capabilities.

## Features
- User registration with role-based access (manager/member)
- JWT-based authentication
- Password hashing with bcrypt
- PostgreSQL database with GORM ORM
- GraphQL API with type safety

## Tech Stack
- **Language**: Go 1.21+
- **GraphQL**: gqlgen
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT tokens
- **Password Hashing**: bcrypt

## Getting Started

### Prerequisites
- Go 1.21 or higher
- PostgreSQL database
- Git

### Installation
1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up your PostgreSQL database
4. Update database connection in `cmd/main.go`
5. Run the service:
   ```bash
   go run cmd/main.go
   ```

The GraphQL playground will be available at `http://localhost:8080/`

## GraphQL Schema

```graphql
scalar UUID

type User {
  id: UUID!
  username: String!
  email: String!
  role: String!
}

type AuthPayload {
  token: String!
  user: User!
}

type Query {
  fetchUsers: [User!]!
}

type Mutation {
  createUser(username: String!, email: String!, password: String!, role: String!): User!
  login(email: String!, password: String!): AuthPayload!
  logout: Boolean!
}
```

## GraphQL Query Examples

### 1. Create User (Mutation)

**Manager User:**
```graphql
mutation CreateManager {
  createUser(
    username: "john_doe"
    email: "john.doe@example.com"
    password: "securePassword123"
    role: "manager"
  ) {
    id
    username
    email
    role
  }
}
```

**Member User:**
```graphql
mutation CreateMember {
  createUser(
    username: "jane_smith"
    email: "jane.smith@example.com"
    password: "anotherSecurePass456"
    role: "member"
  ) {
    id
    username
    email
    role
  }
}
```

**Expected Response:**
```json
{
  "data": {
    "createUser": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "username": "john_doe",
      "email": "john.doe@example.com",
      "role": "manager"
    }
  }
}
```

### 2. Login (Mutation)

```graphql
mutation Login {
  login(
    email: "john.doe@example.com"
    password: "securePassword123"
  ) {
    token
    user {
      id
      username
      email
      role
    }
  }
}
```

**Expected Response:**
```json
{
  "data": {
    "login": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "username": "john_doe",
        "email": "john.doe@example.com",
        "role": "manager"
      }
    }
  }
}
```

### 3. Logout (Mutation)

```graphql
mutation Logout {
  logout
}
```

**Expected Response:**
```json
{
  "data": {
    "logout": true
  }
}
```

### 4. Fetch All Users (Query)

```graphql
query FetchAllUsers {
  fetchUsers {
    id
    username
    email
    role
  }
}
```

**Expected Response:**
```json
{
  "data": {
    "fetchUsers": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "username": "john_doe",
        "email": "john.doe@example.com",
        "role": "manager"
      },
      {
        "id": "987fcdeb-51a2-43d6-b789-123456789abc",
        "username": "jane_smith",
        "email": "jane.smith@example.com",
        "role": "member"
      }
    ]
  }
}
```

## Authentication Headers

For protected endpoints (when implemented), include the JWT token in the Authorization header:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Error Handling

### Common Error Responses

**Invalid Role:**
```json
{
  "errors": [
    {
      "message": "invalid role",
      "path": ["createUser"]
    }
  ]
}
```

**Invalid Credentials:**
```json
{
  "errors": [
    {
      "message": "invalid credentials",
      "path": ["login"]
    }
  ]
}
```

**Email Already Exists:**
```json
{
  "errors": [
    {
      "message": "UNIQUE constraint failed: users.email",
      "path": ["createUser"]
    }
  ]
}
```

## Business Rules

### User Roles
- **manager**: Can create teams, add/remove members, manage assets
- **member**: Can participate in teams, manage personal assets

### Validation Rules
- Email must be unique
- Role must be either "manager" or "member"
- Password is hashed using bcrypt with cost 12
- All fields are required for user creation

## Database Schema

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    role VARCHAR(10) NOT NULL,
    password_hash VARCHAR NOT NULL
);
```

## Service Architecture

### Repository Layer
- `Repository` interface defines data access methods
- `GormRepository` implements the interface using GORM
- Methods: `Create`, `FindByEmail`, `FetchAll`, `Login`

### Service Layer
- `Service` struct handles business logic
- Methods: `CreateUser`, `Login`, `GetUserByEmail`, `FetchUsers`
- Includes password hashing and role validation

### GraphQL Layer
- Resolvers handle GraphQL operations
- Located in `graph/user.resolvers.go`
- Auto-generated types in `graph/model/models_gen.go`

## Testing with GraphQL Playground

1. Start the service: `go run cmd/main.go`
2. Open your browser to `http://localhost:8080/`
3. Use the GraphQL Playground to test queries and mutations
4. Copy and paste the examples above to test the API
