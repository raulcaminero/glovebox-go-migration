import React from 'react';
import { ExternalLink, GitBranch } from 'lucide-react';

export const Header: React.FC = () => {
  return (
    <header className="brand-header">
      <div className="brand">
        <div className="brand-logo">GB</div>
        <div className="brand-title">
          <h1>GloveBox CRM Migration Control Center</h1>
          <p>Modern React + Vite Frontend ➔ Node/NestJS & Go Strangler Architecture</p>
        </div>
      </div>
      <div className="status-cluster">
        <div className="status-pill">
          <span className="status-dot"></span>
          Gateway :8000
        </div>
        <div className="status-pill">
          <span className="status-dot"></span>
          Go Service :8080
        </div>
        <div className="status-pill">
          <span className="status-dot amber"></span>
          NestJS :3001
        </div>
        <div className="status-pill">
          <span className="status-dot"></span>
          Postgres :5435
        </div>
        <a
          href="https://github.com/raulcaminero/glovebox-go-migration"
          target="_blank"
          rel="noreferrer"
          className="status-pill github-pill"
        >
          <GitBranch size={14} />
          GitHub Repo
          <ExternalLink size={12} />
        </a>
      </div>
    </header>
  );
};
