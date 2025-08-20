import React, { useState } from 'react';
import { useTeams, useAuth, useUsers, useFolders } from '../hooks/useApi';
import { teamService } from '../services/teamService';
import { assetService } from '../services/assetService';
import type { User, Team, TeamMember } from '../types';

const Teams: React.FC = () => {
  const { teams, loading, error, createTeam, refetch } = useTeams();
  const { users, loading: usersLoading } = useUsers();
  const { folders } = useFolders();
  const { isManager, user } = useAuth();
  
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newTeamName, setNewTeamName] = useState('');
  const [showAddMemberForm, setShowAddMemberForm] = useState<string | null>(null);
  const [showAddManagerForm, setShowAddManagerForm] = useState<string | null>(null);
  const [showShareAssetsForm, setShowShareAssetsForm] = useState<string | null>(null);
  const [selectedUserId, setSelectedUserId] = useState('');
  const [selectedFolderId, setSelectedFolderId] = useState('');

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

  const handleAddMember = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedUserId || !showAddMemberForm) return;
    
    try {
      await teamService.addMember(showAddMemberForm, { userId: selectedUserId });
      setSelectedUserId('');
      setShowAddMemberForm(null);
      refetch(); // Refresh teams data
    } catch (err) {
      console.error('Error adding member:', err);
    }
  };

  const handleAddManager = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedUserId || !showAddManagerForm) return;
    
    try {
      await teamService.addManager(showAddManagerForm, { userId: selectedUserId });
      setSelectedUserId('');
      setShowAddManagerForm(null);
      refetch(); // Refresh teams data
    } catch (err) {
      console.error('Error adding manager:', err);
    }
  };

  const handleRemoveMember = async (teamId: string, memberId: string) => {
    if (window.confirm('Are you sure you want to remove this member?')) {
      try {
        await teamService.removeMember(teamId, memberId);
        refetch(); // Refresh teams data
      } catch (err) {
        console.error('Error removing member:', err);
      }
    }
  };

  const handleRemoveManager = async (teamId: string, managerId: string) => {
    if (window.confirm('Are you sure you want to remove this manager?')) {
      try {
        await teamService.removeManager(teamId, managerId);
        refetch(); // Refresh teams data
      } catch (err) {
        console.error('Error removing manager:', err);
      }
    }
  };

  // Helper function to get available users for adding to team
  const getAvailableUsers = (team: Team | undefined, role: 'member' | 'manager') => {
    if (!users || !team) return [];
    
    const existingUserIds = [
      ...(team.members?.map((m: TeamMember) => m.userId) || []),
      ...(team.managers?.map((m: TeamMember) => m.userId) || [])
    ];
    
    return users.filter((user: User) => {
      // Don't show users already in the team
      if (existingUserIds.includes(user.id)) return false;
      
      // For manager role, only show users with manager role
      if (role === 'manager') return user.role === 'manager';
      
      // For member role, show all remaining users
      return true;
    });
  };

  // Helper function to check if current user can manage team
  const canManageTeam = (team: Team) => {
    if (!user || !isManager) return false;
    
    // Check if current user is a manager of this team
    return team.managers?.some((manager: TeamMember) => manager.userId === user.id) || false;
  };

  const handleShareAssetsWithTeam = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!showShareAssetsForm || !selectedFolderId) return;

    const team = teams?.find(t => t.id === showShareAssetsForm);
    if (!team) return;

    try {
      // Share folder with all team members
      const allMembers = [...(team.members || []), ...(team.managers || [])];
      
      for (const member of allMembers) {
        if (member.userId !== user?.id) { // Don't share with self
          await assetService.shareFolder(selectedFolderId, {
            userId: member.userId,
            permission: 'read'
          });
        }
      }

      setShowShareAssetsForm(null);
      setSelectedFolderId('');
      alert('Assets shared with team successfully!');
    } catch (err) {
      console.error('Error sharing assets with team:', err);
      alert('Failed to share assets with team');
    }
  };

  if (loading || usersLoading) return <div className="loading">Loading teams...</div>;
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

      {/* Add Member Form */}
      {showAddMemberForm && (
        <div className="create-form">
          <form onSubmit={handleAddMember}>
            <h3>Add Member to Team</h3>
            <div className="form-group">
              <label htmlFor="memberSelect">Select User:</label>
              <select
                id="memberSelect"
                value={selectedUserId}
                onChange={(e) => setSelectedUserId(e.target.value)}
                required
              >
                <option value="">Choose a user...</option>
                {getAvailableUsers(teams?.find(t => t.id === showAddMemberForm), 'member')?.map((user) => (
                  <option key={user.id} value={user.id}>
                    {user.username} ({user.email}) - {user.role}
                  </option>
                ))}
              </select>
            </div>
            <div className="form-buttons">
              <button type="submit">Add Member</button>
              <button type="button" onClick={() => {
                setShowAddMemberForm(null);
                setSelectedUserId('');
              }}>
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Add Manager Form */}
      {showAddManagerForm && (
        <div className="create-form">
          <form onSubmit={handleAddManager}>
            <h3>Add Manager to Team</h3>
            <div className="form-group">
              <label htmlFor="managerSelect">Select Manager:</label>
              <select
                id="managerSelect"
                value={selectedUserId}
                onChange={(e) => setSelectedUserId(e.target.value)}
                required
              >
                <option value="">Choose a manager...</option>
                {getAvailableUsers(teams?.find(t => t.id === showAddManagerForm), 'manager')?.map((user) => (
                  <option key={user.id} value={user.id}>
                    {user.username} ({user.email})
                  </option>
                ))}
              </select>
            </div>
            <div className="form-buttons">
              <button type="submit">Add Manager</button>
              <button type="button" onClick={() => {
                setShowAddManagerForm(null);
                setSelectedUserId('');
              }}>
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Share Assets with Team Form */}
      {showShareAssetsForm && (
        <div className="create-form">
          <form onSubmit={handleShareAssetsWithTeam}>
            <h3>Share Assets with Team</h3>
            <div className="form-group">
              <label htmlFor="folderSelect">Select Folder to Share:</label>
              <select
                id="folderSelect"
                value={selectedFolderId}
                onChange={(e) => setSelectedFolderId(e.target.value)}
                required
              >
                <option value="">Choose a folder...</option>
                {folders?.map((folder) => (
                  <option key={folder.id} value={folder.id}>
                    {folder.folderName} ({folder.notes?.length || 0} notes)
                  </option>
                ))}
              </select>
            </div>
            <div className="form-buttons">
              <button type="submit">Share with Team</button>
              <button type="button" onClick={() => {
                setShowShareAssetsForm(null);
                setSelectedFolderId('');
              }}>
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="teams-list">
        {teams?.length === 0 ? (
          <p>No teams found. {isManager ? 'Create your first team!' : 'No teams available.'}</p>
        ) : (
          teams?.map((team) => (
            <div key={team.id} className="team-card">
              <h3>{team.teamName}</h3>
              <div className="team-info">
                <div className="team-members">
                  <div className="team-section-header">
                    <h4>Managers ({team.managers?.length || 0})</h4>
                    {canManageTeam(team) && (
                      <button 
                        className="add-btn"
                        onClick={() => setShowAddManagerForm(team.id)}
                      >
                        Add Manager
                      </button>
                    )}
                  </div>
                  {team.managers && team.managers.length > 0 && (
                    <ul>
                      {team.managers?.map((manager) => (
                        <li key={manager.userId} className="member-item">
                          <span>{manager.userName} ({manager.email})</span>
                          {canManageTeam(team) && (
                            <button 
                              className="remove-btn"
                              onClick={() => handleRemoveManager(team.id, manager.userId)}
                            >
                              Remove
                            </button>
                          )}
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
                <div className="team-members">
                  <div className="team-section-header">
                    <h4>Members ({team.members?.length || 0})</h4>
                    {canManageTeam(team) && (
                      <button 
                        className="add-btn"
                        onClick={() => setShowAddMemberForm(team.id)}
                      >
                        Add Member
                      </button>
                    )}
                  </div>
                  {team.members && team.members.length > 0 && (
                    <ul>
                      {team.members?.map((member) => (
                        <li key={member.userId} className="member-item">
                          <span>{member.userName} ({member.email})</span>
                          {canManageTeam(team) && (
                            <button 
                              className="remove-btn"
                              onClick={() => handleRemoveMember(team.id, member.userId)}
                            >
                              Remove
                            </button>
                          )}
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
              {canManageTeam(team) && (
                <div className="team-actions">
                  <button 
                    className="share-assets-btn"
                    onClick={() => setShowShareAssetsForm(team.id)}
                  >
                    Share Assets with Team
                  </button>
                </div>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default Teams;
