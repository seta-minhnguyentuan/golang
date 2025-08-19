import React, { useState, useEffect } from 'react';
import { assetService } from '../services/assetService';
import type { Folder, Note } from '../types';

const Assets: React.FC = () => {
  const [folders, setFolders] = useState<Folder[]>([]);
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [activeTab, setActiveTab] = useState<'folders' | 'notes'>('folders');
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newFolderName, setNewFolderName] = useState('');
  const [newNote, setNewNote] = useState({
    title: '',
    content: '',
    folder_id: '',
  });

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [foldersData, notesData] = await Promise.all([
        assetService.getFolders(),
        assetService.getNotes(),
      ]);
      setFolders(foldersData);
      setNotes(notesData);
    } catch (err) {
      setError('Failed to load assets');
      console.error('Error loading assets:', err);
    } finally {
      setLoading(false);
    }
  };

  const createFolder = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await assetService.createFolder({ name: newFolderName });
      setNewFolderName('');
      setShowCreateForm(false);
      loadData();
    } catch (err) {
      setError('Failed to create folder');
      console.error('Error creating folder:', err);
    }
  };

  const createNote = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await assetService.createNote(newNote);
      setNewNote({ title: '', content: '', folder_id: '' });
      setShowCreateForm(false);
      loadData();
    } catch (err) {
      setError('Failed to create note');
      console.error('Error creating note:', err);
    }
  };

  const deleteFolder = async (folderId: string) => {
    if (window.confirm('Are you sure you want to delete this folder?')) {
      try {
        await assetService.deleteFolder(folderId);
        loadData();
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
        loadData();
      } catch (err) {
        setError('Failed to delete note');
        console.error('Error deleting note:', err);
      }
    }
  };

  if (loading) return <div>Loading assets...</div>;
  if (error) return <div className="error">{error}</div>;

  return (
    <div className="assets-container">
      <div className="assets-header">
        <h2>Asset Management</h2>
        <div className="tabs">
          <button
            className={activeTab === 'folders' ? 'active' : ''}
            onClick={() => setActiveTab('folders')}
          >
            Folders
          </button>
          <button
            className={activeTab === 'notes' ? 'active' : ''}
            onClick={() => setActiveTab('notes')}
          >
            Notes
          </button>
        </div>
        <button onClick={() => setShowCreateForm(true)}>
          Create {activeTab === 'folders' ? 'Folder' : 'Note'}
        </button>
      </div>

      {showCreateForm && (
        <div className="create-form">
          {activeTab === 'folders' ? (
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
                />
              </div>
              <div className="form-buttons">
                <button type="submit">Create</button>
                <button type="button" onClick={() => setShowCreateForm(false)}>
                  Cancel
                </button>
              </div>
            </form>
          ) : (
            <form onSubmit={createNote}>
              <h3>Create New Note</h3>
              <div className="form-group">
                <label htmlFor="noteTitle">Title:</label>
                <input
                  type="text"
                  id="noteTitle"
                  value={newNote.title}
                  onChange={(e) => setNewNote({ ...newNote, title: e.target.value })}
                  required
                />
              </div>
              <div className="form-group">
                <label htmlFor="noteFolder">Folder:</label>
                <select
                  id="noteFolder"
                  value={newNote.folder_id}
                  onChange={(e) => setNewNote({ ...newNote, folder_id: e.target.value })}
                  required
                >
                  <option value="">Select a folder</option>
                  {folders.map((folder) => (
                    <option key={folder.id} value={folder.id}>
                      {folder.name}
                    </option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label htmlFor="noteContent">Content:</label>
                <textarea
                  id="noteContent"
                  value={newNote.content}
                  onChange={(e) => setNewNote({ ...newNote, content: e.target.value })}
                  rows={5}
                  required
                />
              </div>
              <div className="form-buttons">
                <button type="submit">Create</button>
                <button type="button" onClick={() => setShowCreateForm(false)}>
                  Cancel
                </button>
              </div>
            </form>
          )}
        </div>
      )}

      <div className="assets-content">
        {activeTab === 'folders' ? (
          <div className="folders-list">
            {folders.length === 0 ? (
              <p>No folders found. Create your first folder!</p>
            ) : (
              folders.map((folder) => (
                <div key={folder.id} className="folder-card">
                  <h3>{folder.name}</h3>
                  <div className="folder-dates">
                    <small>Created: {new Date(folder.created_at).toLocaleDateString()}</small>
                  </div>
                  <div className="folder-actions">
                    <button onClick={() => deleteFolder(folder.id)}>Delete</button>
                  </div>
                </div>
              ))
            )}
          </div>
        ) : (
          <div className="notes-list">
            {notes.length === 0 ? (
              <p>No notes found. Create your first note!</p>
            ) : (
              notes.map((note) => (
                <div key={note.id} className="note-card">
                  <h3>{note.title}</h3>
                  <p>{note.content.substring(0, 100)}...</p>
                  <div className="note-info">
                    <small>Folder ID: {note.folder_id}</small>
                    <small>Created: {new Date(note.created_at).toLocaleDateString()}</small>
                  </div>
                  <div className="note-actions">
                    <button onClick={() => deleteNote(note.id)}>Delete</button>
                  </div>
                </div>
              ))
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default Assets;
