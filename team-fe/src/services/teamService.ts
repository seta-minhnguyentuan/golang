import { userApi } from './api';
import type { Team, CreateTeamRequest, AddMemberRequest, TeamsResponse } from '../types';

export const teamService = {
  async getAllTeams(): Promise<Team[]> {
    const response = await userApi.get<TeamsResponse>('/teams');
    return response.data.teams;
  },

  async createTeam(teamData: CreateTeamRequest): Promise<Team> {
    const response = await userApi.post<Team>('/teams', teamData);
    return response.data;
  },

  async getTeam(teamId: string): Promise<Team> {
    const response = await userApi.get<Team>(`/teams/${teamId}`);
    return response.data;
  },

  async addMember(teamId: string, memberData: AddMemberRequest): Promise<{ message: string }> {
    const response = await userApi.post<{ message: string }>(`/teams/${teamId}/members`, memberData);
    return response.data;
  },

  async removeMember(teamId: string, memberId: string): Promise<{ message: string }> {
    const response = await userApi.delete<{ message: string }>(`/teams/${teamId}/members/${memberId}`);
    return response.data;
  },

  async addManager(teamId: string, managerData: AddMemberRequest): Promise<{ message: string }> {
    const response = await userApi.post<{ message: string }>(`/teams/${teamId}/managers`, managerData);
    return response.data;
  },

  async removeManager(teamId: string, managerId: string): Promise<{ message: string }> {
    const response = await userApi.delete<{ message: string }>(`/teams/${teamId}/managers/${managerId}`);
    return response.data;
  },
};
