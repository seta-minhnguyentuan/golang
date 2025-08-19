import React from 'react';
import { useAuth } from '../contexts/AuthContext';

const Header: React.FC = () => {
  const { user, logout, isAuthenticated } = useAuth();

  const handleLogout = async () => {
    try {
      await logout();
    } catch (err) {
      console.error('Logout error:', err);
    }
  };

  if (!isAuthenticated) return null;

  return (
    <header className="header">
      <div className="header-content">
        <h1>Team Management App</h1>
        <div className="user-info">
          <span>Welcome, {user?.username}!</span>
          <span className="user-role">({user?.role})</span>
          <button onClick={handleLogout} className="logout-btn">
            Logout
          </button>
        </div>
      </div>
    </header>
  );
};

export default Header;
