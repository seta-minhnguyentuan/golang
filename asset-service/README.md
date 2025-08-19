# Asset Service

A RESTful API service for managing folders and notes, built with Go and Gin framework.

## Features

- **Folder Management**: Create, list, retrieve, and delete folders
- **Note Management**: Create notes within folders
- **Health Check**: Service health monitoring endpoint

## Getting Started

### Prerequisites

- Go 1.19 or higher
- PostgreSQL database
- Required environment variables (see Configuration section)

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up your database and environment variables
4. Run the service:
   ```bash
   go run cmd/api/main.go
   ```

### Configuration

The service uses the following environment variables:

- `HTTP_PORT`: Server port (default: 8080)
- Database configuration variables (see `internal/config/database.go`)

## API Documentation

Base URL: `http://localhost:8080`

### Health Check

#### Check Service Health
```bash
curl -X GET http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok"
}
```

### Folder APIs

#### Create a Folder
```bash
curl -X POST http://localhost:8080/api/v1/folders \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Documents"
  }'
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "folderName": "My Documents",
  "notes": [],
  "sharings": [],
  "createdAt": "2025-08-13T10:00:00Z",
  "updatedAt": "2025-08-13T10:00:00Z"
}
```

#### List All Folders
```bash
curl -X GET http://localhost:8080/api/v1/folders
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "folderName": "My Documents",
    "notes": [],
    "sharings": [],
    "createdAt": "2025-08-13T10:00:00Z",
    "updatedAt": "2025-08-13T10:00:00Z"
  }
]
```

#### Get Folder by ID
```bash
curl -X GET http://localhost:8080/api/v1/folders/550e8400-e29b-41d4-a716-446655440000
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "folderName": "My Documents",
  "notes": [
    {
      "id": "660f9500-f30c-52e5-b827-557766551111",
      "noteName": "Meeting Notes",
      "noteContent": "Important meeting discussions",
      "folderId": "550e8400-e29b-41d4-a716-446655440000",
      "sharings": [],
      "createdAt": "2025-08-13T10:30:00Z",
      "updatedAt": "2025-08-13T10:30:00Z"
    }
  ],
  "sharings": [],
  "createdAt": "2025-08-13T10:00:00Z",
  "updatedAt": "2025-08-13T10:00:00Z"
}
```

#### Delete a Folder
```bash
curl -X DELETE http://localhost:8080/api/v1/folders/550e8400-e29b-41d4-a716-446655440000
```

**Response:** Status `204 No Content`

### Note APIs

#### Create a Note
```bash
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Meeting Notes",
    "content": "Important meeting discussions and action items",
    "folder_id": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

**Response:**
```json
{
  "id": "660f9500-f30c-52e5-b827-557766551111",
  "noteName": "Meeting Notes",
  "noteContent": "Important meeting discussions and action items",
  "folderId": "550e8400-e29b-41d4-a716-446655440000",
  "sharings": [],
  "createdAt": "2025-08-13T10:30:00Z",
  "updatedAt": "2025-08-13T10:30:00Z"
}
```

## Error Responses

The API returns standard HTTP status codes and error messages in JSON format:

### Bad Request (400)
```json
{
  "error": "validation error message"
}
```

### Internal Server Error (500)
```json
{
  "error": "internal server error message"
}
```

## Data Models

### Folder
- `id`: UUID (auto-generated)
- `folderName`: String (required)
- `notes`: Array of Note objects
- `sharings`: Array of Sharing objects
- `createdAt`: Timestamp
- `updatedAt`: Timestamp

### Note
- `id`: UUID (auto-generated)
- `noteName`: String
- `noteContent`: String
- `folderId`: UUID (required, foreign key)
- `sharings`: Array of Sharing objects
- `createdAt`: Timestamp
- `updatedAt`: Timestamp

# Asset Service Sharing API Documentation

## Overview

The Asset Service Sharing API allows users to share folders and notes with other users, providing either read or write access. Only the owner of an asset can share it or revoke sharing permissions.

## Authentication

All sharing endpoints require authentication via JWT token. The token should be included in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## API Endpoints

### Folder Sharing

#### Share a Folder
```
POST /api/v1/folders/{folderId}/share
```

**Request Body:**
```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "permission": "read" // or "write"
}
```

**Response:**
```json
{
  "message": "Folder shared successfully"
}
```

#### Revoke Folder Sharing
```
DELETE /api/v1/folders/{folderId}/share/{userId}
```

**Response:**
```json
{
  "message": "Folder sharing revoked successfully"
}
```

#### List Folder Sharings
```
GET /api/v1/folders/{folderId}/share
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "userId": "550e8400-e29b-41d4-a716-446655440000",
    "permission": "read",
    "folderId": "550e8400-e29b-41d4-a716-446655440002",
    "createdAt": "2025-08-19T10:30:00Z",
    "updatedAt": "2025-08-19T10:30:00Z"
  }
]
```

### Note Sharing

#### Share a Note
```
POST /api/v1/notes/{noteId}/share
```

**Request Body:**
```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "permission": "read" // or "write"
}
```

**Response:**
```json
{
  "message": "Note shared successfully"
}
```

#### Revoke Note Sharing
```
DELETE /api/v1/notes/{noteId}/share/{userId}
```

**Response:**
```json
{
  "message": "Note sharing revoked successfully"
}
```

#### List Note Sharings
```
GET /api/v1/notes/{noteId}/share
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "userId": "550e8400-e29b-41d4-a716-446655440000",
    "permission": "write",
    "noteId": "550e8400-e29b-41d4-a716-446655440004",
    "createdAt": "2025-08-19T10:30:00Z",
    "updatedAt": "2025-08-19T10:30:00Z"
  }
]
```

## Permission Types

- **read**: User can view the shared asset but cannot modify it
- **write**: User can view and modify the shared asset

## Business Rules

1. **Owner Only**: Only the owner of a folder or note can share it or manage sharing permissions
2. **Self-Sharing Prevention**: Users cannot share assets with themselves
3. **Folder Inheritance**: When a folder is shared, all notes within that folder are implicitly shared with the same permissions
4. **Permission Updates**: If a user already has access to an asset, sharing it again will update their permission level
5. **Manager Access**: Managers can view (read-only) all assets their team members own or have access to (this functionality is planned for future implementation)

## Error Responses

All endpoints return appropriate HTTP status codes and error messages:

```json
{
  "error": "Error description"
}
```

Common error scenarios:
- **400 Bad Request**: Invalid UUID format, invalid permission type, business rule violations
- **401 Unauthorized**: Missing or invalid authentication token
- **404 Not Found**: Asset not found
- **500 Internal Server Error**: Server-side errors

## Example Usage

### Sharing a Folder with Read Access

```bash
curl -X POST \
  http://localhost:8080/api/v1/folders/550e8400-e29b-41d4-a716-446655440002/share \
  -H 'Authorization: Bearer your-jwt-token' \
  -H 'Content-Type: application/json' \
  -d '{
    "userId": "550e8400-e29b-41d4-a716-446655440000",
    "permission": "read"
  }'
```

### Upgrading Permission to Write Access

```bash
curl -X POST \
  http://localhost:8080/api/v1/folders/550e8400-e29b-41d4-a716-446655440002/share \
  -H 'Authorization: Bearer your-jwt-token' \
  -H 'Content-Type: application/json' \
  -d '{
    "userId": "550e8400-e29b-41d4-a716-446655440000",
    "permission": "write"
  }'
```

### Revoking Access

```bash
curl -X DELETE \
  http://localhost:8080/api/v1/folders/550e8400-e29b-41d4-a716-446655440002/share/550e8400-e29b-41d4-a716-446655440000 \
  -H 'Authorization: Bearer your-jwt-token'
```

## Implementation Notes

The sharing API is implemented with:

- **Repository Layer**: Handles database operations for sharing records
- **Service Layer**: Implements business logic and permission validation
- **Handler Layer**: Handles HTTP requests and responses
- **Models**: Defines FolderSharing and NoteSharing entities with UUID primary keys

The implementation ensures data consistency and proper validation of all sharing operations while maintaining security through ownership verification and authentication requirements.
