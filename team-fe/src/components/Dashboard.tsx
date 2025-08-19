import React, { useState } from 'react';
import Teams from './Teams';
import Assets from './Assets';

const Dashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'teams' | 'assets'>('teams');

  return (
    <div className="dashboard">
      <nav className="dashboard-nav">
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

      <div className="dashboard-content">
        {activeTab === 'teams' ? <Teams /> : <Assets />}
      </div>
    </div>
  );
};

export default Dashboard;
