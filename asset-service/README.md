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

