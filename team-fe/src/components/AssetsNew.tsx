import React, { useState } from 'react';
import { useFolders, useFolder, useTeams, useSharing } from '../hooks/useApi';

const Assets: React.FC = () => {
  const { folders, loading: foldersLoading, createFolder, deleteFolder } = useFolders();
  const { teams } = useTeams();
  const { shareTeamAssets } = useSharing();
  
  const [selectedFolderId, setSelectedFolderId] = useState<string | null>(null);
  const [showCreateFolder, setShowCreateFolder] = useState(false);
  const [showCreateNote, setShowCreateNote] = useState(false);
  const [showSharing, setShowSharing] = useState(false);
  const [newFolderName, setNewFolderName] = useState('');
  const [newNoteTitle, setNewNoteTitle] = useState('');
  const [newNoteContent, setNewNoteContent] = useState('');

  const { folder, loading: folderLoading, createNote } = useFolder(selectedFolderId);

  const handleCreateFolder = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newFolderName.trim()) return;
    
    try {
      await createFolder(newFolderName);
      setNewFolderName('');
      setShowCreateFolder(false);
    } catch (err) {
      console.error('Error creating folder:', err);
    }
  };

  const handleCreateNote = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newNoteTitle.trim() || !newNoteContent.trim()) return;
    
    try {
      await createNote(newNoteTitle, newNoteContent);
      setNewNoteTitle('');
      setNewNoteContent('');
      setShowCreateNote(false);
    } catch (err) {
      console.error('Error creating note:', err);
    }
  };

  const handleDeleteFolder = async (folderId: string) => {
    if (confirm('Are you sure you want to delete this folder and all its notes?')) {
      try {
        await deleteFolder(folderId);
        if (selectedFolderId === folderId) {
          setSelectedFolderId(null);
        }
      } catch (err) {
        console.error('Error deleting folder:', err);
      }
    }
  };

  const handleShareWithTeam = async (teamId: string) => {
    if (!selectedFolderId) return;
    
    try {
      await shareTeamAssets(teamId, selectedFolderId, 'read');
      alert('Folder shared with team members successfully!');
      setShowSharing(false);
    } catch (err) {
      console.error('Error sharing with team:', err);
      alert('Failed to share folder with team');
    }
  };

  if (foldersLoading) return <div className="loading">Loading assets...</div>;

  return (
    <div className="assets-container">
      <div className="assets-header">
        <h2>My Assets</h2>
        <button 
          className="create-btn" 
          onClick={() => setShowCreateFolder(true)}
        >
          Create Folder
        </button>
      </div>

      <div className="assets-content">
        {/* Folders List */}
        <div className="folders-panel">
          <h3>Folders</h3>
          
          {showCreateFolder && (
            <div className="create-form">
              <form onSubmit={handleCreateFolder}>
                <input
                  type="text"
                  value={newFolderName}
                  onChange={(e) => setNewFolderName(e.target.value)}
                  placeholder="Folder name"
                  required
                />
                <button type="submit">Create</button>
                <button type="button" onClick={() => setShowCreateFolder(false)}>
                  Cancel
                </button>
              </form>
            </div>
          )}

          <div className="folders-list">
            {folders.length === 0 ? (
              <p>No folders found. Create your first folder!</p>
            ) : (
              folders.map((folder) => (
                <div 
                  key={folder.id} 
                  className={`folder-item ${selectedFolderId === folder.id ? 'selected' : ''}`}
                  onClick={() => setSelectedFolderId(folder.id)}
                >
                  <h4>{folder.folderName}</h4>
                  <p>{folder.notes?.length || 0} notes</p>
                  <div className="folder-actions">
                    <button onClick={(e) => {
                      e.stopPropagation();
                      handleDeleteFolder(folder.id);
                    }}>
                      Delete
                    </button>
                    {teams?.length > 0 && (
                      <button onClick={(e) => {
                        e.stopPropagation();
                        setSelectedFolderId(folder.id);
                        setShowSharing(true);
                      }}>
                        Share
                      </button>
                    )}
                  </div>
                  <small>Created: {new Date(folder.createdAt).toLocaleDateString()}</small>
                </div>
              ))
            )}
          </div>
        </div>

        {/* Notes Panel */}
        <div className="notes-panel">
          {selectedFolderId ? (
            <>
              <div className="notes-header">
                <h3>Notes in {folder?.folderName}</h3>
                <button onClick={() => setShowCreateNote(true)}>
                  Add Note
                </button>
              </div>

              {showCreateNote && (
                <div className="create-form">
                  <form onSubmit={handleCreateNote}>
                    <input
                      type="text"
                      value={newNoteTitle}
                      onChange={(e) => setNewNoteTitle(e.target.value)}
                      placeholder="Note title"
                      required
                    />
                    <textarea
                      value={newNoteContent}
                      onChange={(e) => setNewNoteContent(e.target.value)}
                      placeholder="Note content"
                      rows={4}
                      required
                    />
                    <div className="form-buttons">
                      <button type="submit">Create Note</button>
                      <button type="button" onClick={() => setShowCreateNote(false)}>
                        Cancel
                      </button>
                    </div>
                  </form>
                </div>
              )}

              {folderLoading ? (
                <div className="loading">Loading notes...</div>
              ) : (
                <div className="notes-list">
                  {folder?.notes.length === 0 ? (
                    <p>No notes in this folder. Add your first note!</p>
                  ) : (
                    folder?.notes.map((note) => (
                      <div key={note.id} className="note-item">
                        <h4>{note.noteName}</h4>
                        <p>{note.noteContent}</p>
                        <small>Created: {new Date(note.createdAt).toLocaleDateString()}</small>
                      </div>
                    ))
                  )}
                </div>
              )}
            </>
          ) : (
            <div className="no-selection">
              <p>Select a folder to view its notes</p>
            </div>
          )}
        </div>
      </div>

      {/* Sharing Modal */}
      {showSharing && (
        <div className="modal-overlay" onClick={() => setShowSharing(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h3>Share Folder with Team</h3>
            <div className="teams-list">
              {teams.map((team) => (
                <div key={team.id} className="team-share-item">
                  <h4>{team.teamName}</h4>
                  <p>{team.members.length + team.managers.length} members</p>
                  <button onClick={() => handleShareWithTeam(team.id)}>
                    Share with Team
                  </button>
                </div>
              ))}
            </div>
            <button className="close-btn" onClick={() => setShowSharing(false)}>
              Close
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default Assets;
