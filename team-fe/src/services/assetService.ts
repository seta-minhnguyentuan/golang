import { assetApi } from './api';
import type { 
  Folder, 
  CreateFolderRequest, 
  Note, 
  CreateNoteRequest, 
  UpdateNoteRequest,
  Sharing,
  ShareRequest
} from '../types';

export const assetService = {
  // Health Check
  async checkHealth(): Promise<{ status: string }> {
    const response = await assetApi.get<{ status: string }>('/health');
    return response.data;
  },

  // Folder operations
  async getFolders(): Promise<Folder[]> {
    const response = await assetApi.get<Folder[]>('/folders');
    return response.data;
  },

  async createFolder(folderData: CreateFolderRequest): Promise<Folder> {
    const response = await assetApi.post<Folder>('/folders', folderData);
    return response.data;
  },

  async getFolder(folderId: string): Promise<Folder> {
    const response = await assetApi.get<Folder>(`/folders/${folderId}`);
    return response.data;
  },

  async deleteFolder(folderId: string): Promise<void> {
    await assetApi.delete(`/folders/${folderId}`);
  },

  // Note operations
  async getNotes(): Promise<Note[]> {
    const response = await assetApi.get<Note[]>('/notes');
    return response.data;
  },

  async createNote(noteData: CreateNoteRequest): Promise<Note> {
    const response = await assetApi.post<Note>('/notes', noteData);
    return response.data;
  },

  async getNote(noteId: string): Promise<Note> {
    const response = await assetApi.get<Note>(`/notes/${noteId}`);
    return response.data;
  },

  async updateNote(noteId: string, noteData: UpdateNoteRequest): Promise<Note> {
    const response = await assetApi.put<Note>(`/notes/${noteId}`, noteData);
    return response.data;
  },

  async deleteNote(noteId: string): Promise<void> {
    await assetApi.delete(`/notes/${noteId}`);
  },

  // Folder Sharing operations
  async shareFolder(folderId: string, shareData: ShareRequest): Promise<{ message: string }> {
    const response = await assetApi.post<{ message: string }>(`/folders/${folderId}/share`, shareData);
    return response.data;
  },

  async revokeFolderSharing(folderId: string, userId: string): Promise<{ message: string }> {
    const response = await assetApi.delete<{ message: string }>(`/folders/${folderId}/share/${userId}`);
    return response.data;
  },

  async getFolderSharings(folderId: string): Promise<Sharing[]> {
    const response = await assetApi.get<Sharing[]>(`/folders/${folderId}/share`);
    return response.data;
  },

  // Note Sharing operations
  async shareNote(noteId: string, shareData: ShareRequest): Promise<{ message: string }> {
    const response = await assetApi.post<{ message: string }>(`/notes/${noteId}/share`, shareData);
    return response.data;
  },

  async revokeNoteSharing(noteId: string, userId: string): Promise<{ message: string }> {
    const response = await assetApi.delete<{ message: string }>(`/notes/${noteId}/share/${userId}`);
    return response.data;
  },

  async getNoteSharings(noteId: string): Promise<Sharing[]> {
    const response = await assetApi.get<Sharing[]>(`/notes/${noteId}/share`);
    return response.data;
  },
};
