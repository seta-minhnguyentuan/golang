import React, { useState, useEffect } from 'react';
import { teamService } from '../services/teamService';
import type { Team } from '../types';

const Teams: React.FC = () => {
  const [teams, setTeams] = useState<Team[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newTeamName, setNewTeamName] = useState('');
  const [newTeamDescription, setNewTeamDescription] = useState('');

  useEffect(() => {
    loadTeams();
  }, []);

  const loadTeams = async () => {
    try {
      setLoading(true);
      const teamsData = await teamService.getAllTeams();
      setTeams(teamsData);
    } catch (err) {
      setError('Failed to load teams');
      console.error('Error loading teams:', err);
    } finally {
      setLoading(false);
    }
  };

  const createTeam = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await teamService.createTeam({
        name: newTeamName,
        description: newTeamDescription,
      });
      setNewTeamName('');
      setNewTeamDescription('');
      setShowCreateForm(false);
      loadTeams();
    } catch (err) {
      setError('Failed to create team');
      console.error('Error creating team:', err);
    }
  };

  if (loading) return <div>Loading teams...</div>;
  if (error) return <div className="error">{error}</div>;

  return (
    <div className="teams-container">
      <div className="teams-header">
        <h2>Teams</h2>
        <button onClick={() => setShowCreateForm(true)}>Create Team</button>
      </div>

      {showCreateForm && (
        <div className="create-form">
          <form onSubmit={createTeam}>
            <h3>Create New Team</h3>
            <div className="form-group">
              <label htmlFor="teamName">Team Name:</label>
              <input
                type="text"
                id="teamName"
                value={newTeamName}
                onChange={(e) => setNewTeamName(e.target.value)}
                required
              />
            </div>
            <div className="form-group">
              <label htmlFor="teamDescription">Description:</label>
              <textarea
                id="teamDescription"
                value={newTeamDescription}
                onChange={(e) => setNewTeamDescription(e.target.value)}
              />
            </div>
            <div className="form-buttons">
              <button type="submit">Create</button>
              <button type="button" onClick={() => setShowCreateForm(false)}>
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="teams-list">
        {teams.length === 0 ? (
          <p>No teams found. Create your first team!</p>
        ) : (
          teams.map((team) => (
            <div key={team.id} className="team-card">
              <h3>{team.name}</h3>
              <p>{team.description}</p>
              <div className="team-info">
                <span>Members: {team.members?.length || 0}</span>
                <span>Managers: {team.managers?.length || 0}</span>
              </div>
              <div className="team-dates">
                <small>Created: {new Date(team.created_at).toLocaleDateString()}</small>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default Teams;
