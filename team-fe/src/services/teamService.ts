import { userApi } from './api';
import type { Team, CreateTeamRequest, AddMemberRequest } from '../types';

export const teamService = {
  async getAllTeams(): Promise<Team[]> {
    const response = await userApi.get('/teams');
    return response.data.teams;
  },

  async createTeam(teamData: CreateTeamRequest): Promise<Team> {
    const response = await userApi.post('/teams', teamData);
    return response.data;
  },

  async getTeam(teamId: string): Promise<Team> {
    const response = await userApi.get(`/teams/${teamId}`);
    return response.data;
  },

  async addMember(teamId: string, memberData: AddMemberRequest): Promise<void> {
    await userApi.post(`/teams/${teamId}/members`, memberData);
  },

  async removeMember(teamId: string, memberId: string): Promise<void> {
    await userApi.delete(`/teams/${teamId}/members/${memberId}`);
  },

  async addManager(teamId: string, managerData: AddMemberRequest): Promise<void> {
    await userApi.post(`/teams/${teamId}/managers`, managerData);
  },

  async removeManager(teamId: string, managerId: string): Promise<void> {
    await userApi.delete(`/teams/${teamId}/managers/${managerId}`);
  },
};
