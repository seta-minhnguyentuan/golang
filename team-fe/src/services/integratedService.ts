import { userService } from './userService';
import { teamService } from './teamService';
import { assetService } from './assetService';
import type { 
  User, 
  Team, 
  Folder, 
  Note, 
  CreateUserRequest, 
  CreateTeamRequest, 
  CreateFolderRequest, 
  CreateNoteRequest,
  LoginRequest
} from '../types';

/**
 * Integrated service that combines User Service, Team Service, and Asset Service
 * This provides a unified interface for all API operations
 */
export class IntegratedService {
  // User Management (GraphQL)
  static async createUser(userData: CreateUserRequest): Promise<User> {
    return await userService.createUser(userData);
  }

  static async login(credentials: LoginRequest) {
    return await userService.login(credentials);
  }

  static async logout() {
    return await userService.logout();
  }

  static async fetchUsers(): Promise<User[]> {
    return await userService.fetchUsers();
  }

  // Team Management (REST)
  static async getAllTeams(): Promise<Team[]> {
    return await teamService.getAllTeams();
  }

  static async createTeam(teamData: CreateTeamRequest): Promise<Team> {
    return await teamService.createTeam(teamData);
  }

  static async getTeam(teamId: string): Promise<Team> {
    return await teamService.getTeam(teamId);
  }

  static async addTeamMember(teamId: string, userId: string) {
    return await teamService.addMember(teamId, { userId });
  }

  static async removeTeamMember(teamId: string, memberId: string) {
    return await teamService.removeMember(teamId, memberId);
  }

  static async addTeamManager(teamId: string, userId: string) {
    return await teamService.addManager(teamId, { userId });
  }

  static async removeTeamManager(teamId: string, managerId: string) {
    return await teamService.removeManager(teamId, managerId);
  }

  // Asset Management (REST)
  static async checkAssetServiceHealth() {
    return await assetService.checkHealth();
  }

  static async getFolders(): Promise<Folder[]> {
    return await assetService.getFolders();
  }

  static async createFolder(folderData: CreateFolderRequest): Promise<Folder> {
    return await assetService.createFolder(folderData);
  }

  static async getFolder(folderId: string): Promise<Folder> {
    return await assetService.getFolder(folderId);
  }

  static async deleteFolder(folderId: string): Promise<void> {
    return await assetService.deleteFolder(folderId);
  }

  static async getNotes(): Promise<Note[]> {
    return await assetService.getNotes();
  }

  static async createNote(noteData: CreateNoteRequest): Promise<Note> {
    return await assetService.createNote(noteData);
  }

  static async getNote(noteId: string): Promise<Note> {
    return await assetService.getNote(noteId);
  }

  static async updateNote(noteId: string, title: string, content: string): Promise<Note> {
    return await assetService.updateNote(noteId, { title, content });
  }

  static async deleteNote(noteId: string): Promise<void> {
    return await assetService.deleteNote(noteId);
  }

  // Sharing Operations
  static async shareFolder(folderId: string, userId: string, permission: 'read' | 'write' = 'read') {
    return await assetService.shareFolder(folderId, { userId, permission });
  }

  static async revokeFolderSharing(folderId: string, userId: string) {
    return await assetService.revokeFolderSharing(folderId, userId);
  }

  static async getFolderSharings(folderId: string) {
    return await assetService.getFolderSharings(folderId);
  }

  static async shareNote(noteId: string, userId: string, permission: 'read' | 'write' = 'read') {
    return await assetService.shareNote(noteId, { userId, permission });
  }

  static async revokeNoteSharing(noteId: string, userId: string) {
    return await assetService.revokeNoteSharing(noteId, userId);
  }

  static async getNoteSharings(noteId: string) {
    return await assetService.getNoteSharings(noteId);
  }

  // Utility Methods
  static getCurrentUser(): User | null {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  }

  static getToken(): string | null {
    return localStorage.getItem('token');
  }

  static isAuthenticated(): boolean {
    return !!this.getToken();
  }

  static isManager(): boolean {
    const user = this.getCurrentUser();
    return user?.role === 'manager';
  }

  static isMember(): boolean {
    const user = this.getCurrentUser();
    return user?.role === 'member';
  }

  /**
   * Complete workflow example: Create a team with assets
   */
  static async createTeamWithAssets(
    teamName: string, 
    folderName: string,
    initialNote?: { title: string; content: string }
  ) {
    try {
      // 1. Create team
      const team = await this.createTeam({ teamName });

      // 2. Create folder for the team
      const folder = await this.createFolder({ name: `${folderName} - ${teamName}` });

      // 3. Create initial note if provided
      let note;
      if (initialNote) {
        note = await this.createNote({
          title: initialNote.title,
          content: initialNote.content,
          folder_id: folder.id
        });
      }

      return { team, folder, note };
    } catch (error) {
      console.error('Error creating team with assets:', error);
      throw error;
    }
  }

  /**
   * Share team assets with all team members
   */
  static async shareTeamAssets(teamId: string, folderId: string, permission: 'read' | 'write' = 'read') {
    try {
      const team = await this.getTeam(teamId);
      const sharePromises = [];

      // Share with all team members
      for (const member of team.members) {
        sharePromises.push(this.shareFolder(folderId, member.userId, permission));
      }

      // Share with all team managers (with write permission)
      for (const manager of team.managers) {
        sharePromises.push(this.shareFolder(folderId, manager.userId, 'write'));
      }

      await Promise.all(sharePromises);
      return { message: `Assets shared with ${team.members.length + team.managers.length} team members` };
    } catch (error) {
      console.error('Error sharing team assets:', error);
      throw error;
    }
  }
}

export default IntegratedService;
