# Frontend API Integration Guide

This guide explains how to use the integrated APIs for User Service, Team Management, and Asset Service in your React frontend.

## Overview

The frontend now provides a unified interface to interact with:
- **User Service** (GraphQL + REST): User authentication and team management
- **Asset Service** (REST): Folder and note management with sharing capabilities

## API Service Configuration

### Base URLs
```typescript
// src/services/api.ts
export const USER_SERVICE_URL = 'http://localhost:8080';        // User & Team Service
export const ASSET_SERVICE_URL = 'http://localhost:8080/api/v1'; // Asset Service
```

### Authentication
JWT tokens are automatically included in requests for protected endpoints:
- **All team management endpoints** (`/teams`)
- **All asset endpoints** (`/folders`, `/notes`, `/*/share`) - Asset service requires authentication for all operations

## Available Services

### 1. IntegratedService (src/services/integratedService.ts)

A unified service class that combines all API operations:

```typescript
import { IntegratedService } from '../services/integratedService';

// User Management
await IntegratedService.createUser({ username, email, password, role });
await IntegratedService.login({ email, password });
await IntegratedService.logout();
await IntegratedService.fetchUsers();

// Team Management
await IntegratedService.createTeam({ teamName });
await IntegratedService.getAllTeams();
await IntegratedService.addTeamMember(teamId, userId);

// Asset Management
await IntegratedService.createFolder({ name });
await IntegratedService.getFolders();
await IntegratedService.createNote({ title, content, folder_id });
await IntegratedService.shareFolder(folderId, userId, 'read');
```

### 2. React Hooks (src/hooks/useApi.ts)

React hooks for state management and API calls:

#### useAuth Hook
```typescript
const { user, login, logout, createUser, isAuthenticated, isManager } = useAuth();

// Login
const handleLogin = async () => {
  try {
    await login('user@example.com', 'password');
  } catch (error) {
    console.error('Login failed:', error);
  }
};
```

#### useTeams Hook
```typescript
const { teams, loading, error, createTeam, refetch } = useTeams();

// Create team
const handleCreateTeam = async () => {
  try {
    await createTeam('Development Team');
  } catch (error) {
    console.error('Team creation failed:', error);
  }
};
```

#### useFolders Hook
```typescript
const { folders, loading, createFolder, deleteFolder } = useFolders();

// Create folder
const handleCreateFolder = async () => {
  try {
    await createFolder('Project Documents');
  } catch (error) {
    console.error('Folder creation failed:', error);
  }
};
```

#### useFolder Hook
```typescript
const { folder, loading, createNote } = useFolder(selectedFolderId);

// Create note in folder
const handleCreateNote = async () => {
  try {
    await createNote('Meeting Notes', 'Important discussion points...');
  } catch (error) {
    console.error('Note creation failed:', error);
  }
};
```

#### useSharing Hook
```typescript
const { shareFolder, shareTeamAssets } = useSharing();

// Share folder with specific user
await shareFolder(folderId, userId, 'read');

// Share folder with all team members
await shareTeamAssets(teamId, folderId, 'read');
```

## Component Examples

### 1. Teams Component (src/components/Teams.tsx)

Updated to use the new API integration:

```typescript
import React, { useState } from 'react';
import { useTeams, useAuth } from '../hooks/useApi';

const Teams: React.FC = () => {
  const { teams, loading, error, createTeam } = useTeams();
  const { isManager } = useAuth();
  const [newTeamName, setNewTeamName] = useState('');

  const handleCreateTeam = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newTeamName.trim()) return;
    
    try {
      await createTeam(newTeamName);
      setNewTeamName('');
    } catch (err) {
      console.error('Error creating team:', err);
    }
  };

  if (loading) return <div>Loading teams...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h2>Teams</h2>
      {isManager && (
        <form onSubmit={handleCreateTeam}>
          <input
            value={newTeamName}
            onChange={(e) => setNewTeamName(e.target.value)}
            placeholder="Team name"
            required
          />
          <button type="submit">Create Team</button>
        </form>
      )}
      
      {teams.map((team) => (
        <div key={team.id}>
          <h3>{team.teamName}</h3>
          <p>Managers: {team.managers.length}</p>
          <p>Members: {team.members.length}</p>
        </div>
      ))}
    </div>
  );
};
```

### 2. Assets Component (src/components/AssetsNew.tsx)

Complete asset management with sharing:

```typescript
import React, { useState } from 'react';
import { useFolders, useFolder, useTeams, useSharing } from '../hooks/useApi';

const Assets: React.FC = () => {
  const { folders, createFolder, deleteFolder } = useFolders();
  const { teams } = useTeams();
  const { shareTeamAssets } = useSharing();
  const [selectedFolderId, setSelectedFolderId] = useState<string | null>(null);
  
  const { folder, createNote } = useFolder(selectedFolderId);

  // Create folder, notes, share with teams...
  // See full implementation in AssetsNew.tsx
};
```

## API Payload Structures

### User Service Payloads

#### Create User
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securePassword123",
  "role": "manager"
}
```

#### Login Response
```json
{
  "data": {
    "login": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user": {
        "id": "uuid",
        "username": "john_doe",
        "email": "john@example.com",
        "role": "manager"
      }
    }
  }
}
```

#### Team Creation
```json
{
  "teamName": "Development Team",
  "managers": [
    {
      "userId": "uuid",
      "userName": "john_doe"
    }
  ],
  "members": [
    {
      "userId": "uuid", 
      "userName": "jane_smith"
    }
  ]
}
```

### Asset Service Payloads

#### Create Folder
```json
{
  "name": "Project Documents"
}
```

#### Folder Response
```json
{
  "id": "uuid",
  "folderName": "Project Documents",
  "notes": [],
  "sharings": [],
  "createdAt": "2025-08-19T10:00:00Z",
  "updatedAt": "2025-08-19T10:00:00Z"
}
```

#### Create Note
```json
{
  "title": "Meeting Notes",
  "content": "Important discussion points...",
  "folder_id": "folder-uuid"
}
```

#### Share Request
```json
{
  "userId": "user-uuid",
  "permission": "read"
}
```

## Workflow Examples

### Complete Team Setup with Assets

```typescript
const setupTeamWithAssets = async () => {
  try {
    // 1. Create team
    const team = await IntegratedService.createTeam({
      teamName: 'Development Team'
    });

    // 2. Create shared folder
    const folder = await IntegratedService.createFolder({
      name: 'Team Documents'
    });

    // 3. Create initial notes
    await IntegratedService.createNote({
      title: 'Team Guidelines',
      content: 'Team collaboration guidelines and best practices...',
      folder_id: folder.id
    });

    // 4. Share folder with all team members
    await IntegratedService.shareTeamAssets(team.id, folder.id, 'write');

    console.log('Team setup completed successfully!');
  } catch (error) {
    console.error('Team setup failed:', error);
  }
};
```

### User Authentication Flow

```typescript
const AuthFlow = () => {
  const { user, login, logout, isAuthenticated } = useAuth();

  const handleLogin = async (email: string, password: string) => {
    try {
      await login(email, password);
      // User is now authenticated, token stored automatically
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  const handleLogout = async () => {
    await logout();
    // User logged out, token removed automatically
  };

  return (
    <div>
      {isAuthenticated ? (
        <div>
          <p>Welcome, {user?.username}!</p>
          <button onClick={handleLogout}>Logout</button>
        </div>
      ) : (
        <LoginForm onLogin={handleLogin} />
      )}
    </div>
  );
};
```

## Error Handling

All API calls include proper error handling:

```typescript
try {
  const result = await IntegratedService.createTeam({ teamName: 'New Team' });
  console.log('Success:', result);
} catch (error) {
  if (error.response?.status === 401) {
    console.error('Unauthorized - please login');
  } else if (error.response?.status === 400) {
    console.error('Validation error:', error.response.data.error);
  } else {
    console.error('Unexpected error:', error.message);
  }
}
```

## TypeScript Types

All API responses are properly typed. Key types include:

```typescript
interface User {
  id: string;
  username: string;
  email: string;
  role: 'manager' | 'member';
}

interface Team {
  id: string;
  teamName: string;
  managers: TeamMember[];
  members: TeamMember[];
  createdAt: string;
  updatedAt: string;
}

interface Folder {
  id: string;
  folderName: string;
  notes: Note[];
  sharings: Sharing[];
  createdAt: string;
  updatedAt: string;
}

interface Note {
  id: string;
  noteName: string;
  noteContent: string;
  folderId: string;
  sharings: Sharing[];
  createdAt: string;
  updatedAt: string;
}
```

## Getting Started

1. **Install dependencies** (already included in package.json):
   ```bash
   npm install
   ```

2. **Start your backend services**:
   - User Service: `go run cmd/main.go` (port 8080)
   - Asset Service: `go run cmd/api/main.go` (port 8080 or configured)

3. **Start the frontend**:
   ```bash
   npm run dev
   ```

4. **Use the integrated APIs** in your components:
   ```typescript
   import { useAuth, useTeams, useFolders } from '../hooks/useApi';
   import { IntegratedService } from '../services/integratedService';
   ```

The frontend now provides a complete, type-safe interface to interact with all your backend services!
