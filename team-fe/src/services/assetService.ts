import { assetApi } from './api';
import type { 
  Folder, 
  CreateFolderRequest, 
  Note, 
  CreateNoteRequest, 
  UpdateNoteRequest 
} from '../types';

export const assetService = {
  // Folder operations
  async getFolders(): Promise<Folder[]> {
    const response = await assetApi.get('/folders');
    return response.data;
  },

  async createFolder(folderData: CreateFolderRequest): Promise<Folder> {
    const response = await assetApi.post('/folders', folderData);
    return response.data;
  },

  async getFolder(folderId: string): Promise<Folder> {
    const response = await assetApi.get(`/folders/${folderId}`);
    return response.data;
  },

  async deleteFolder(folderId: string): Promise<void> {
    await assetApi.delete(`/folders/${folderId}`);
  },

  // Note operations
  async getNotes(): Promise<Note[]> {
    const response = await assetApi.get('/notes');
    return response.data;
  },

  async createNote(noteData: CreateNoteRequest): Promise<Note> {
    const response = await assetApi.post('/notes', noteData);
    return response.data;
  },

  async getNote(noteId: string): Promise<Note> {
    const response = await assetApi.get(`/notes/${noteId}`);
    return response.data;
  },

  async updateNote(noteId: string, noteData: UpdateNoteRequest): Promise<Note> {
    const response = await assetApi.put(`/notes/${noteId}`, noteData);
    return response.data;
  },

  async deleteNote(noteId: string): Promise<void> {
    await assetApi.delete(`/notes/${noteId}`);
  },

  // Sharing operations
  async shareFolder(folderId: string, userId: string): Promise<void> {
    await assetApi.post(`/folders/${folderId}/share`, { user_id: userId });
  },

  async revokeFolderSharing(folderId: string, userId: string): Promise<void> {
    await assetApi.delete(`/folders/${folderId}/share/${userId}`);
  },

  async getFolderSharings(folderId: string): Promise<unknown[]> {
    const response = await assetApi.get(`/folders/${folderId}/share`);
    return response.data;
  },

  async shareNote(noteId: string, userId: string): Promise<void> {
    await assetApi.post(`/notes/${noteId}/share`, { user_id: userId });
  },

  async revokeNoteSharing(noteId: string, userId: string): Promise<void> {
    await assetApi.delete(`/notes/${noteId}/share/${userId}`);
  },

  async getNoteSharings(noteId: string): Promise<unknown[]> {
    const response = await assetApi.get(`/notes/${noteId}/share`);
    return response.data;
  },
};
