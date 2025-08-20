import React, { useState, useEffect } from 'react';
import { assetService } from '../services/assetService';
import { useUsers } from '../hooks/useApi';
import type { Folder, User } from '../types';

const AssetsEnhanced: React.FC = () => {
  const [folders, setFolders] = useState<Folder[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateFolderForm, setShowCreateFolderForm] = useState(false);
  const [showCreateNoteForm, setShowCreateNoteForm] = useState<string | null>(null);
  const [showShareForm, setShowShareForm] = useState<{ type: 'folder' | 'note', id: string } | null>(null);
  const [expandedFolders, setExpandedFolders] = useState<Set<string>>(new Set());
  const [newFolderName, setNewFolderName] = useState('');
  const [newNote, setNewNote] = useState({
    title: '',
    content: '',
  });
  const [shareData, setShareData] = useState({
    userId: '',
    permission: 'read' as 'read' | 'write'
  });

  const { users, loading: usersLoading } = useUsers();

  useEffect(() => {
    loadFolders();
  }, []);

  const loadFolders = async () => {
    try {
      setLoading(true);
      const foldersData = await assetService.getFolders();
      setFolders(foldersData);
    } catch (err) {
      setError('Failed to load folders');
      console.error('Error loading folders:', err);
    } finally {
      setLoading(false);
    }
  };

  const createFolder = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await assetService.createFolder({ name: newFolderName });
      setNewFolderName('');
      setShowCreateFolderForm(false);
      loadFolders();
    } catch (err) {
      setError('Failed to create folder');
      console.error('Error creating folder:', err);
    }
  };

  const createNote = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!showCreateNoteForm) return;
    
    try {
      await assetService.createNote({
        title: newNote.title,
        content: newNote.content,
        folder_id: showCreateNoteForm
      });
      setNewNote({ title: '', content: '' });
      setShowCreateNoteForm(null);
      loadFolders(); // Reload to get updated folder with new note
    } catch (err) {
      setError('Failed to create note');
      console.error('Error creating note:', err);
    }
  };

  const deleteFolder = async (folderId: string) => {
    if (window.confirm('Are you sure you want to delete this folder and all its notes?')) {
      try {
        await assetService.deleteFolder(folderId);
        loadFolders();
      } catch (err) {
        setError('Failed to delete folder');
        console.error('Error deleting folder:', err);
      }
    }
  };

  const deleteNote = async (noteId: string) => {
    if (window.confirm('Are you sure you want to delete this note?')) {
      try {
        await assetService.deleteNote(noteId);
        loadFolders(); // Reload to update folder data
      } catch (err) {
        setError('Failed to delete note');
        console.error('Error deleting note:', err);
      }
    }
  };

  const toggleFolder = (folderId: string) => {
    const newExpanded = new Set(expandedFolders);
    if (newExpanded.has(folderId)) {
      newExpanded.delete(folderId);
    } else {
      newExpanded.add(folderId);
    }
    setExpandedFolders(newExpanded);
  };

  const handleShare = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!showShareForm || !shareData.userId) return;

    try {
      if (showShareForm.type === 'folder') {
        await assetService.shareFolder(showShareForm.id, {
          userId: shareData.userId,
          permission: shareData.permission
        });
      } else {
        await assetService.shareNote(showShareForm.id, {
          userId: shareData.userId,
          permission: shareData.permission
        });
      }
      
      setShowShareForm(null);
      setShareData({ userId: '', permission: 'read' });
      loadFolders(); // Refresh to show updated sharing info
    } catch (err) {
      setError(`Failed to share ${showShareForm.type}`);
      console.error(`Error sharing ${showShareForm.type}:`, err);
    }
  };

  const handleRevokeSharing = async (type: 'folder' | 'note', itemId: string, userId: string) => {
    if (!window.confirm('Are you sure you want to revoke sharing?')) return;

    try {
      if (type === 'folder') {
        await assetService.revokeFolderSharing(itemId, userId);
      } else {
        await assetService.revokeNoteSharing(itemId, userId);
      }
      loadFolders(); // Refresh to show updated sharing info
    } catch (err) {
      setError(`Failed to revoke ${type} sharing`);
      console.error(`Error revoking ${type} sharing:`, err);
    }
  };

  if (loading) return <div className="loading">Loading assets...</div>;
  if (error) return <div className="error">{error}</div>;

  return (
    <div className="assets-enhanced-container">
      <div className="assets-header">
        <h2>Asset Management</h2>
        <button 
          className="create-btn"
          onClick={() => setShowCreateFolderForm(true)}
        >
          Create Folder
        </button>
      </div>

      {/* Create Folder Form */}
      {showCreateFolderForm && (
        <div className="create-form">
          <form onSubmit={createFolder}>
            <h3>Create New Folder</h3>
            <div className="form-group">
              <label htmlFor="folderName">Folder Name:</label>
              <input
                type="text"
                id="folderName"
                value={newFolderName}
                onChange={(e) => setNewFolderName(e.target.value)}
                required
                placeholder="Enter folder name"
              />
            </div>
            <div className="form-buttons">
              <button type="submit">Create</button>
              <button type="button" onClick={() => setShowCreateFolderForm(false)}>
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Create Note Form */}
      {showCreateNoteForm && (
        <div className="create-form">
          <form onSubmit={createNote}>
            <h3>Add Note to Folder</h3>
            <div className="form-group">
              <label htmlFor="noteTitle">Note Title:</label>
              <input
                type="text"
                id="noteTitle"
                value={newNote.title}
                onChange={(e) => setNewNote({ ...newNote, title: e.target.value })}
                required
                placeholder="Enter note title"
              />
            </div>
            <div className="form-group">
              <label htmlFor="noteContent">Note Content:</label>
              <textarea
                id="noteContent"
                value={newNote.content}
                onChange={(e) => setNewNote({ ...newNote, content: e.target.value })}
                rows={5}
                required
                placeholder="Enter note content"
              />
            </div>
            <div className="form-buttons">
              <button type="submit">Add Note</button>
              <button type="button" onClick={() => {
                setShowCreateNoteForm(null);
                setNewNote({ title: '', content: '' });
              }}>
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Share Form */}
      {showShareForm && (
        <div className="create-form">
          <form onSubmit={handleShare}>
            <h3>Share {showShareForm.type === 'folder' ? 'Folder' : 'Note'}</h3>
            <div className="form-group">
              <label htmlFor="userSelect">Select User:</label>
              <select
                id="userSelect"
                value={shareData.userId}
                onChange={(e) => setShareData({ ...shareData, userId: e.target.value })}
                required
              >
                <option value="">Choose a user...</option>
                {users?.map((user: User) => (
                  <option key={user.id} value={user.id}>
                    {user.username} ({user.email}) - {user.role}
                  </option>
                ))}
              </select>
            </div>
            <div className="form-group">
              <label htmlFor="permissionSelect">Permission:</label>
              <select
                id="permissionSelect"
                value={shareData.permission}
                onChange={(e) => setShareData({ ...shareData, permission: e.target.value as 'read' | 'write' })}
              >
                <option value="read">Read Only</option>
                <option value="write">Read & Write</option>
              </select>
            </div>
            <div className="form-buttons">
              <button type="submit" disabled={usersLoading}>
                Share
              </button>
              <button type="button" onClick={() => {
                setShowShareForm(null);
                setShareData({ userId: '', permission: 'read' });
              }}>
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Folders List */}
      <div className="folders-list">
        {folders?.length === 0 ? (
          <div className="empty-state">
            <p>No folders found.</p>
            <p>Create your first folder to organize your notes!</p>
          </div>
        ) : (
          folders?.map((folder) => (
            <div key={folder.id} className="folder-card">
              <div className="folder-header">
                <div className="folder-info">
                  <button 
                    className="folder-toggle"
                    onClick={() => toggleFolder(folder.id)}
                    aria-label={expandedFolders.has(folder.id) ? 'Collapse folder' : 'Expand folder'}
                  >
                    {expandedFolders.has(folder.id) ? '📂' : '📁'}
                  </button>
                  <h3 onClick={() => toggleFolder(folder.id)} style={{ cursor: 'pointer' }}>
                    {folder.folderName}
                  </h3>
                  <span className="notes-count">({folder.notes?.length || 0} notes)</span>
                </div>
                <div className="folder-actions">
                  <button 
                    className="add-note-btn"
                    onClick={() => setShowCreateNoteForm(folder.id)}
                  >
                    Add Note
                  </button>
                  <button 
                    className="share-btn"
                    onClick={() => setShowShareForm({ type: 'folder', id: folder.id })}
                  >
                    Share
                  </button>
                  <button 
                    className="delete-btn"
                    onClick={() => deleteFolder(folder.id)}
                  >
                    Delete Folder
                  </button>
                </div>
              </div>

              <div className="folder-meta">
                <small>Created: {new Date(folder.createdAt).toLocaleDateString()}</small>
                {folder.sharings && folder.sharings.length > 0 && (
                  <small>Shared with {folder.sharings.length} user(s)</small>
                )}
              </div>

              {/* Sharing Information */}
              {folder.sharings && folder.sharings.length > 0 && (
                <div className="sharing-info">
                  <h5>Shared with:</h5>
                  <ul className="sharing-list">
                    {folder.sharings.map((sharing) => (
                      <li key={sharing.id} className="sharing-item">
                        <span>User ID: {sharing.userId} ({sharing.permission})</span>
                        <button 
                          className="revoke-btn"
                          onClick={() => handleRevokeSharing('folder', folder.id, sharing.userId)}
                        >
                          Revoke
                        </button>
                      </li>
                    ))}
                  </ul>
                </div>
              )}

              {/* Notes List - Only show when folder is expanded */}
              {expandedFolders.has(folder.id) && (
                <div className="notes-list">
                  {folder.notes && folder.notes.length > 0 ? (
                    <div className="notes-grid">
                      {folder.notes?.map((note) => (
                        <div key={note.id} className="note-card">
                          <div className="note-header">
                            <h4>{note.noteName}</h4>
                            <div className="note-actions">
                              <button 
                                className="share-note-btn"
                                onClick={() => setShowShareForm({ type: 'note', id: note.id })}
                                aria-label="Share note"
                              >
                                📤
                              </button>
                              <button 
                                className="delete-note-btn"
                                onClick={() => deleteNote(note.id)}
                                aria-label="Delete note"
                              >
                                ×
                              </button>
                            </div>
                          </div>
                          <div className="note-content">
                            <p>{note.noteContent?.length > 150 
                              ? note.noteContent.substring(0, 150) + '...' 
                              : note.noteContent}
                            </p>
                          </div>
                          <div className="note-meta">
                            <small>Created: {new Date(note.createdAt).toLocaleDateString()}</small>
                            {note.sharings && note.sharings.length > 0 && (
                              <small>Shared with {note.sharings.length} user(s)</small>
                            )}
                          </div>
                          
                          {/* Note Sharing Information */}
                          {note.sharings && note.sharings.length > 0 && (
                            <div className="note-sharing-info">
                              <h6>Shared with:</h6>
                              <ul className="sharing-list">
                                {note.sharings.map((sharing) => (
                                  <li key={sharing.id} className="sharing-item">
                                    <span>User ID: {sharing.userId} ({sharing.permission})</span>
                                    <button 
                                      className="revoke-btn"
                                      onClick={() => handleRevokeSharing('note', note.id, sharing.userId)}
                                    >
                                      Revoke
                                    </button>
                                  </li>
                                ))}
                              </ul>
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="no-notes">
                      <p>No notes in this folder yet.</p>
                      <button 
                        className="add-first-note-btn"
                        onClick={() => setShowCreateNoteForm(folder.id)}
                      >
                        Add your first note
                      </button>
                    </div>
                  )}
                </div>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default AssetsEnhanced;
