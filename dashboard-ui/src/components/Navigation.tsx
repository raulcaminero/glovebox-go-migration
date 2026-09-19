import React from 'react';
import { GitBranch, Users, Layers, Bot } from 'lucide-react';

interface NavigationProps {
  activeTab: string;
  setActiveTab: (tab: string) => void;
}

export const Navigation: React.FC<NavigationProps> = ({ activeTab, setActiveTab }) => {
  return (
    <div className="nav-tabs">
      <button
        className={`tab-btn ${activeTab === 'strangler' ? 'active' : ''}`}
        onClick={() => setActiveTab('strangler')}
      >
        <GitBranch size={16} /> Strangler Routing Map
      </button>
      <button
        className={`tab-btn ${activeTab === 'crm' ? 'active' : ''}`}
        onClick={() => setActiveTab('crm')}
      >
        <Users size={16} /> CRM Data Manager
      </button>
      <button
        className={`tab-btn ${activeTab === 'adrs' ? 'active' : ''}`}
        onClick={() => setActiveTab('adrs')}
      >
        <Layers size={16} /> Hexagonal Architecture & ADRs
      </button>
      <button
        className={`tab-btn ${activeTab === 'ai' ? 'active' : ''}`}
        onClick={() => setActiveTab('ai')}
      >
        <Bot size={16} /> AI Tools & Architecture Guardrails
      </button>
    </div>
  );
};
