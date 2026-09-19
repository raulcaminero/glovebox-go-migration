import React, { useState } from 'react';
import { Header } from './components/Header';
import { Navigation } from './components/Navigation';
import { StranglerVisualizer } from './components/StranglerVisualizer';
import { HexagonalArch } from './components/HexagonalArch';
import { CrmManager } from './components/CrmManager';
import { AiGuardrails } from './components/AiGuardrails';
import { DocModal } from './components/DocModal';
import { CreateModal } from './components/CreateModal';

export const App: React.FC = () => {
  const [activeTab, setActiveTab] = useState<string>('strangler');
  const [docModalInfo, setDocModalInfo] = useState<{ id: string | null; title: string }>({
    id: null,
    title: '',
  });
  const [isCreateOpen, setIsCreateOpen] = useState<boolean>(false);
  const [refreshKey, setRefreshKey] = useState<number>(0);

  const handleOpenDoc = (docId: string, title: string) => {
    setDocModalInfo({ id: docId, title });
  };

  const handleCloseDoc = () => {
    setDocModalInfo({ id: null, title: '' });
  };

  return (
    <div className="app-container">
      <Header />
      <Navigation activeTab={activeTab} setActiveTab={setActiveTab} />

      {activeTab === 'strangler' && <StranglerVisualizer />}
      {activeTab === 'crm' && (
        <CrmManager key={refreshKey} onOpenCreate={() => setIsCreateOpen(true)} />
      )}
      {activeTab === 'adrs' && <HexagonalArch onOpenDoc={handleOpenDoc} />}
      {activeTab === 'ai' && <AiGuardrails onOpenDoc={handleOpenDoc} />}

      <DocModal
        docId={docModalInfo.id}
        title={docModalInfo.title}
        onClose={handleCloseDoc}
      />

      <CreateModal
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onCreated={() => setRefreshKey((prev) => prev + 1)}
      />
    </div>
  );
};

export default App;
