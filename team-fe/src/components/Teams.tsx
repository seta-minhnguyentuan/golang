import React, { useState } from 'react';
import { useTeams, useAuth } from '../hooks/useApi';

const Teams: React.FC = () => {
  const { teams, loading, error, createTeam } = useTeams();
  const { isManager } = useAuth();
  
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newTeamName, setNewTeamName] = useState('');

  const handleCreateTeam = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newTeamName.trim()) return;
    
    try {
      await createTeam(newTeamName);
      setNewTeamName('');
      setShowCreateForm(false);
    } catch (err) {
      console.error('Error creating team:', err);
    }
  };

  if (loading) return <div className="loading">Loading teams...</div>;
  if (error) return <div className="error">Error: {error}</div>;

  return (
    <div className="teams-container">
      <div className="teams-header">
        <h2>Teams</h2>
        {isManager && (
          <button 
            className="create-btn" 
            onClick={() => setShowCreateForm(true)}
          >
            Create Team
          </button>
        )}
      </div>

      {showCreateForm && (
        <div className="create-form">
          <form onSubmit={handleCreateTeam}>
            <h3>Create New Team</h3>
            <div className="form-group">
              <label htmlFor="teamName">Team Name:</label>
              <input
                type="text"
                id="teamName"
                value={newTeamName}
                onChange={(e) => setNewTeamName(e.target.value)}
                required
                placeholder="Enter team name"
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
          <p>No teams found. {isManager ? 'Create your first team!' : 'No teams available.'}</p>
        ) : (
          teams.map((team) => (
            <div key={team.id} className="team-card">
              <h3>{team.teamName}</h3>
              <div className="team-info">
                <div className="team-members">
                  <h4>Managers ({team.managers?.length || 0})</h4>
                  {team.managers?.length > 0 && (
                    <ul>
                      {team.managers.map((manager) => (
                        <li key={manager.userId}>
                          {manager.userName} ({manager.email})
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
                <div className="team-members">
                  <h4>Members ({team.members?.length || 0})</h4>
                  {team.members?.length > 0 && (
                    <ul>
                      {team.members.map((member) => (
                        <li key={member.userId}>
                          {member.userName} ({member.email})
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
              </div>
              <div className="team-dates">
                <small>Created: {new Date(team.createdAt).toLocaleDateString()}</small>
                <small>Updated: {new Date(team.updatedAt).toLocaleDateString()}</small>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default Teams;
