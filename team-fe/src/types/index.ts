export interface User {
  id: string;
  username: string;
  email: string;
  role: string;
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
  role: string;
}

export interface Team {
  id: string;
  name: string;
  description?: string;
  members: User[];
  managers: User[];
  created_at: string;
  updated_at: string;
}

export interface CreateTeamRequest {
  name: string;
  description?: string;
}

export interface AddMemberRequest {
  user_id: string;
}

export interface Folder {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface CreateFolderRequest {
  name: string;
}

export interface Note {
  id: string;
  title: string;
  content: string;
  folder_id: string;
  created_at: string;
  updated_at: string;
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
