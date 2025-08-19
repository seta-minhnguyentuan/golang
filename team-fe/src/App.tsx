import React, { useState } from 'react';
import { ApolloProvider } from '@apollo/client';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { apolloClient } from './services/userService';
import Login from './components/Login';
import Header from './components/Header';
import Teams from './components/Teams';
import Assets from './components/Assets';
import './App.css';

const AppContent: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const [activeTab, setActiveTab] = useState<'teams' | 'assets'>('teams');

  if (!isAuthenticated) {
    return <Login />;
  }

  return (
    <div className="app">
      <Header />
      <nav className="main-nav">
        <button
          className={activeTab === 'teams' ? 'active' : ''}
          onClick={() => setActiveTab('teams')}
        >
          Teams
        </button>
        <button
          className={activeTab === 'assets' ? 'active' : ''}
          onClick={() => setActiveTab('assets')}
        >
          Assets
        </button>
      </nav>
      <main className="main-content">
        {activeTab === 'teams' ? <Teams /> : <Assets />}
      </main>
    </div>
  );
};

function App() {
  return (
    <ApolloProvider client={apolloClient}>
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </ApolloProvider>
  );
}

export default App;
