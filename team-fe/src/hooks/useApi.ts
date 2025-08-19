import { useState, useEffect, useCallback } from 'react';
import { IntegratedService } from '../services/integratedService';
import type { User, Team, Folder } from '../types';

// Hook for user management
export const useUsers = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchUsers = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await IntegratedService.fetchUsers();
      setUsers(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch users');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  return { users, loading, error, refetch: fetchUsers };
};

// Hook for team management
export const useTeams = () => {
  const [teams, setTeams] = useState<Team[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchTeams = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await IntegratedService.getAllTeams();
      setTeams(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch teams');
    } finally {
      setLoading(false);
    }
  };

  const createTeam = async (teamName: string) => {
    setLoading(true);
    setError(null);
    try {
      const newTeam = await IntegratedService.createTeam({ teamName });
      setTeams(prev => [...prev, newTeam]);
      return newTeam;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create team');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (IntegratedService.isAuthenticated()) {
      fetchTeams();
    }
  }, []);

  return { teams, loading, error, createTeam, refetch: fetchTeams };
};

// Hook for folder management
export const useFolders = () => {
  const [folders, setFolders] = useState<Folder[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchFolders = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await IntegratedService.getFolders();
      setFolders(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch folders');
    } finally {
      setLoading(false);
    }
  };

  const createFolder = async (name: string) => {
    setLoading(true);
    setError(null);
    try {
      const newFolder = await IntegratedService.createFolder({ name });
      setFolders(prev => [...prev, newFolder]);
      return newFolder;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create folder');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const deleteFolder = async (folderId: string) => {
    setLoading(true);
    setError(null);
    try {
      await IntegratedService.deleteFolder(folderId);
      setFolders(prev => prev.filter(folder => folder.id !== folderId));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete folder');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchFolders();
  }, []);

  return { folders, loading, error, createFolder, deleteFolder, refetch: fetchFolders };
};

// Hook for a specific folder with its notes
export const useFolder = (folderId: string | null) => {
  const [folder, setFolder] = useState<Folder | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchFolder = useCallback(async () => {
    if (!folderId) return;
    
    setLoading(true);
    setError(null);
    try {
      const data = await IntegratedService.getFolder(folderId);
      setFolder(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch folder');
    } finally {
      setLoading(false);
    }
  }, [folderId]);

  const createNote = async (title: string, content: string) => {
    if (!folderId) throw new Error('No folder selected');
    
    setLoading(true);
    setError(null);
    try {
      const newNote = await IntegratedService.createNote({
        title,
        content,
        folder_id: folderId
      });
      
      // Update the folder's notes
      setFolder(prev => prev ? {
        ...prev,
        notes: [...prev.notes, newNote]
      } : null);
      
      return newNote;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create note');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchFolder();
  }, [fetchFolder]);

  return { folder, loading, error, createNote, refetch: fetchFolder };
};

// Hook for authentication
export const useAuth = () => {
  const [user, setUser] = useState<User | null>(IntegratedService.getCurrentUser());
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const login = async (email: string, password: string) => {
    setLoading(true);
    setError(null);
    try {
      const authPayload = await IntegratedService.login({ email, password });
      setUser(authPayload.user);
      return authPayload;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const logout = async () => {
    setLoading(true);
    setError(null);
    try {
      await IntegratedService.logout();
      setUser(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Logout failed');
    } finally {
      setLoading(false);
    }
  };

  const createUser = async (username: string, email: string, password: string, role: 'manager' | 'member') => {
    setLoading(true);
    setError(null);
    try {
      const newUser = await IntegratedService.createUser({ username, email, password, role });
      return newUser;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'User creation failed');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return {
    user,
    loading,
    error,
    login,
    logout,
    createUser,
    isAuthenticated: !!user,
    isManager: user?.role === 'manager',
    isMember: user?.role === 'member'
  };
};

// Hook for sharing operations
export const useSharing = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const shareFolder = async (folderId: string, userId: string, permission: 'read' | 'write' = 'read') => {
    setLoading(true);
    setError(null);
    try {
      const result = await IntegratedService.shareFolder(folderId, userId, permission);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to share folder');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const shareNote = async (noteId: string, userId: string, permission: 'read' | 'write' = 'read') => {
    setLoading(true);
    setError(null);
    try {
      const result = await IntegratedService.shareNote(noteId, userId, permission);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to share note');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const shareTeamAssets = async (teamId: string, folderId: string, permission: 'read' | 'write' = 'read') => {
    setLoading(true);
    setError(null);
    try {
      const result = await IntegratedService.shareTeamAssets(teamId, folderId, permission);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to share team assets');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return {
    loading,
    error,
    shareFolder,
    shareNote,
    shareTeamAssets
  };
};
