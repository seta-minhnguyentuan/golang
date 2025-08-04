# User & Team Management Service

## Overview
This is a hybrid API service that provides both GraphQL and REST endpoints for user and team management. It features user authentication, role-based authorization, and comprehensive team management capabilities.

## Features
- **User Management**: Registration, authentication with JWT tokens
- **Team Management**: Create teams, manage members and managers
- **Role-based Access**: Managers can create/manage teams, members can participate
- **Hybrid API**: GraphQL for user operations, REST for team operations
- **Database Relations**: Proper foreign key constraints between users and teams
- **Password Security**: bcrypt hashing with secure defaults

## Tech Stack
- **Language**: Go 1.21+
- **GraphQL**: gqlgen (for user management)
- **REST API**: Gin framework (for team management)
- **Database**: PostgreSQL with UUID support
- **ORM**: GORM with auto-migration
- **Authentication**: JWT tokens with middleware
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
3. Set up your PostgreSQL database:
   ```sql
   CREATE DATABASE user_service;
   CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
   ```
4. Update database connection in `cmd/main.go` if needed
5. Run the service:
   ```bash
   go run cmd/main.go
   ```

## API Endpoints

### GraphQL (User Management)
- **Endpoint**: `POST /user/query`
- **Playground**: `GET /user/query` 
- **Operations**: User registration, login, logout, fetch users

### REST API (Team Management) 
- **Base Path**: `/teams`
- **Authentication**: Requires JWT token in Authorization header
- **Operations**: Create teams, manage members and managers

## API Structure

```
📁 User Management (GraphQL)
└── /user/query
    ├── createUser(username, email, password, role): User
    ├── login(email, password): AuthPayload  
    ├── logout(): Boolean
    └── fetchUsers(): [User]

📁 Team Management (REST)
└── /teams (🔒 JWT Required)
    ├── GET /teams                                    # Get all teams
    ├── POST /teams                                   # Create team
    ├── GET /teams/{teamId}                          # Get team details
    ├── POST /teams/{teamId}/members                 # Add member
    ├── DELETE /teams/{teamId}/members/{memberId}    # Remove member
    ├── POST /teams/{teamId}/managers                # Add manager
    └── DELETE /teams/{teamId}/managers/{managerId}  # Remove manager
```

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

## GraphQL Examples (User Management)

### Access GraphQL Playground
Open your browser to: `http://localhost:8080/user/query`

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
```

## REST API Examples (Team Management)

### 🔑 Authentication Setup
First, get a JWT token by logging in via GraphQL:

```bash
# Login via GraphQL to get JWT token
curl -X POST http://localhost:8080/user/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { login(email: \"john.doe@example.com\", password: \"securePassword123\") { token user { id username role } } }"
  }'
```

**Response:**
```json
{
  "data": {
    "login": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "username": "john_doe",
        "role": "manager"
      }
    }
  }
}
```

### 📋 Team Management Operations

**⚠️ Note**: Replace `YOUR_JWT_TOKEN` with the actual token from login response.

#### 1. Create Team

```bash
curl -X POST http://localhost:8080/teams \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "teamName": "Development Team",
    "managers": [
      {
        "userId": "987fcdeb-51a2-43d6-b789-123456789abc",
        "userName": "jane_manager"
      }
    ],
    "members": [
      {
        "userId": "456e7890-e12b-34d5-a678-901234567def",
        "userName": "bob_developer"
      }
    ]
  }'
```

**Response:**
```json
{
  "id": "team-uuid-here",
  "teamName": "Development Team",
  "managers": [
    {
      "userId": "123e4567-e89b-12d3-a456-426614174000",
      "userName": "john_doe",
      "email": "john.doe@example.com",
      "role": "manager",
      "joinedAt": "2025-08-04T16:30:00Z"
    }
  ],
  "members": [
    {
      "userId": "456e7890-e12b-34d5-a678-901234567def",
      "userName": "bob_developer",
      "email": "bob@example.com",
      "role": "member",
      "joinedAt": "2025-08-04T16:30:00Z"
    }
  ],
  "createdAt": "2025-08-04T16:30:00Z",
  "updatedAt": "2025-08-04T16:30:00Z"
}
```

#### 2. Get All Teams

```bash
curl -X GET http://localhost:8080/teams \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "teams": [
    {
      "id": "team-uuid-1",
      "teamName": "Development Team",
      "managers": [
        {
          "userId": "123e4567-e89b-12d3-a456-426614174000",
          "userName": "john_doe",
          "email": "john.doe@example.com",
          "role": "manager",
          "joinedAt": "2025-08-04T16:30:00Z"
        }
      ],
      "members": [
        {
          "userId": "456e7890-e12b-34d5-a678-901234567def",
          "userName": "bob_developer",
          "email": "bob@example.com",
          "role": "member",
          "joinedAt": "2025-08-04T16:30:00Z"
        }
      ],
      "createdAt": "2025-08-04T16:30:00Z",
      "updatedAt": "2025-08-04T16:30:00Z"
    }
  ]
}
```

**Note**: 
- **Managers** can see all teams in the system
- **Members** can only see teams they belong to

#### 3. Get Team Details

```bash
curl -X GET http://localhost:8080/teams/{team-id} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 4. Add Member to Team

```bash
curl -X POST http://localhost:8080/teams/{team-id}/members \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "userId": "new-member-uuid"
  }'
```

**Response:**
```json
{
  "message": "Member added successfully"
}
```

#### 5. Remove Member from Team

```bash
curl -X DELETE http://localhost:8080/teams/{team-id}/members/{member-id} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "message": "Member removed successfully"
}
```

#### 6. Add Manager to Team

```bash
curl -X POST http://localhost:8080/teams/{team-id}/managers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "userId": "manager-user-uuid"
  }'
```

**Response:**
```json
{
  "message": "Manager added successfully"
}
```

#### 7. Remove Manager from Team

```bash
curl -X DELETE http://localhost:8080/teams/{team-id}/managers/{manager-id} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "message": "Manager removed successfully"
}
```

### 🚫 Error Examples

**Unauthorized Access:**
```json
{
  "error": "Authorization header required"
}
```

**Invalid Token:**
```json
{
  "error": "Invalid token"
}
```

**Permission Denied:**
```json
{
  "error": "only managers can create teams"
}
```

**User Already in Team:**
```json
{
  "error": "user is already a member of this team"
}
```

## Authentication Headers

For team management endpoints, include the JWT token in the Authorization header:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Business Rules

### User Roles
- **manager**: Can create teams, add/remove members and managers, manage assets
- **member**: Can participate in teams, manage personal assets

### Team Management Rules
- Only managers can create teams
- Only team managers can add/remove members and managers
- Users cannot be added to the same team twice
- Team creators are automatically added as managers
- Managers being added to teams must have "manager" role in the system

### Validation Rules
- Email must be unique across all users
- Role must be either "manager" or "member"
- Password is hashed using bcrypt with cost 12
- All fields are required for user creation
- Team names are required for team creation

## Database Schema

### Tables with Relationships

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    role VARCHAR(10) NOT NULL CHECK (role IN ('manager', 'member')),
    password_hash VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Teams table
CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_name VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Team members junction table with foreign key relationships
CREATE TABLE team_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(10) NOT NULL CHECK (role IN ('manager', 'member')),
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Foreign key constraints
    CONSTRAINT fk_team_members_team_id 
        FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    CONSTRAINT fk_team_members_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Unique constraint to prevent duplicate team memberships
    CONSTRAINT idx_team_user UNIQUE (team_id, user_id)
);
```

## Service Architecture

### Domain-Driven Structure
```
├── internal/
│   ├── user/                   # User domain
│   │   ├── handler.go          # HTTP handlers (minimal for GraphQL)
│   │   ├── service.go          # Business logic
│   │   ├── repository.go       # Database logic
│   │   └── model.go            # Structs and DTOs
│   ├── team/                   # Team domain
│   │   ├── handler.go          # REST API handlers
│   │   ├── service.go          # Business logic
│   │   ├── repository.go       # Database logic
│   │   └── model.go            # Structs and DTOs
│   └── auth/                   # Authentication
│       ├── jwt.go              # Token generation/validation
│       └── middleware.go       # JWT middleware
├── router/                     # Route definitions
│   └── router.go               # Central routing configuration
└── graph/                      # GraphQL schema and resolvers
    ├── resolver.go
    ├── user.resolvers.go
    └── schema/
```

### Repository Layer
- **User Repository**: `Create`, `FindByEmail`, `FindByID`, `FetchAll`, `Login`
- **Team Repository**: `Create`, `FindByID`, `FindMembersByTeamID`, `AddMember`, `RemoveMember`
- Interface-based design for easy testing and mocking

### Service Layer
- **User Service**: `CreateUser`, `Login`, `GetUserByEmail`, `FetchUsers`
- **Team Service**: `CreateTeam`, `GetTeamByID`, `AddMember`, `RemoveMember`, `AddManager`, `RemoveManager`
- Includes business logic, validation, and role-based access control

### API Layer
- **GraphQL**: Auto-generated resolvers handle user operations
- **REST**: Gin handlers for team management with JWT middleware
- Centralized routing in `router/router.go`

## Testing Guide

### 1. Start the Service
```bash
go run cmd/main.go
```

### 2. Test GraphQL (User Management)
- Open GraphQL Playground: `http://localhost:8080/user/query`
- Create a manager user first (needed for team operations)
- Login to get JWT token

### 3. Test REST API (Team Management)
- Use curl commands from the examples above
- Ensure you have a valid JWT token from GraphQL login
- Test team creation, member management, and error scenarios

### 4. Database Verification
Connect to your PostgreSQL database to verify relationships:

```sql
-- Check all users
SELECT * FROM users;

-- Check teams with their members
SELECT 
    t.team_name,
    u.username,
    tm.role,
    tm.joined_at
FROM teams t
JOIN team_members tm ON t.id = tm.team_id
JOIN users u ON tm.user_id = u.id
ORDER BY t.team_name, tm.role;
```

## Complete Workflow Example

### Step 1: Create Users (GraphQL)
```bash
# Create a manager
curl -X POST http://localhost:8080/user/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { createUser(username: \"john_manager\", email: \"john@example.com\", password: \"pass123\", role: \"manager\") { id username role } }"
  }'

# Create a member  
curl -X POST http://localhost:8080/user/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { createUser(username: \"jane_member\", email: \"jane@example.com\", password: \"pass123\", role: \"member\") { id username role } }"
  }'
```

### Step 2: Login (GraphQL)
```bash
curl -X POST http://localhost:8080/user/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { login(email: \"john@example.com\", password: \"pass123\") { token user { id } } }"
  }'
```

### Step 3: Create Team (REST)
```bash
curl -X POST http://localhost:8080/teams \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "teamName": "My First Team",
    "members": []
  }'
```

### Step 4: Add Member to Team (REST)
```bash
curl -X POST http://localhost:8080/teams/{team-id}/members \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "userId": "jane-member-uuid"
  }'
```

The service now provides a complete user and team management system with proper database relationships, role-based access control, and both GraphQL and REST API endpoints!
