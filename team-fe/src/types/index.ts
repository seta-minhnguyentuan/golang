// User Service Types
export interface User {
  id: string;
  username: string;
  email: string;
  role: 'manager' | 'member';
}

export interface AuthPayload {
  token: string;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
  role: 'manager' | 'member';
}

// Team Management Types
export interface TeamMember {
  userId: string;
  userName: string;
  email: string;
  role: 'manager' | 'member';
  joinedAt: string;
}

export interface Team {
  id: string;
  teamName: string;
  managers: TeamMember[];
  members: TeamMember[];
  createdAt: string;
  updatedAt: string;
}

export interface CreateTeamRequest {
  teamName: string;
  managers?: {
    userId: string;
    userName: string;
  }[];
  members?: {
    userId: string;
    userName: string;
  }[];
}

export interface AddMemberRequest {
  userId: string;
}

export interface TeamsResponse {
  teams: Team[];
}

// Asset Service Types
export interface Folder {
  id: string;
  folderName: string;
  notes: Note[];
  sharings: Sharing[];
  createdAt: string;
  updatedAt: string;
}

export interface CreateFolderRequest {
  name: string;
}

export interface Note {
  id: string;
  noteName: string;
  noteContent: string;
  folderId: string;
  sharings: Sharing[];
  createdAt: string;
  updatedAt: string;
}

export interface CreateNoteRequest {
  title: string;
  content: string;
  folder_id: string;
}

export interface UpdateNoteRequest {
  title: string;
  content: string;
}

// Sharing Types
export interface Sharing {
  id: string;
  userId: string;
  permission: 'read' | 'write';
  folderId?: string;
  noteId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ShareRequest {
  userId: string;
  permission: 'read' | 'write';
}

// API Response Types
export interface ApiResponse<T> {
  data: T;
  message?: string;
}

export interface ErrorResponse {
  error: string;
}
