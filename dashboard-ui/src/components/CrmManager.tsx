import React, { useState, useEffect } from 'react';
import type { Policyholder } from '../types';

const jwtToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMyJ9.ah_KrNhrpocJD5ZKfyjoKMdCZMst4lvLSayT5B_d8Bk";

interface CrmManagerProps {
  onOpenCreate: () => void;
}

export const CrmManager: React.FC<CrmManagerProps> = ({ onOpenCreate }) => {
  const [list, setList] = useState<Policyholder[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const loadPolicyholders = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch('/api/v1/policyholders', {
        headers: { Authorization: `Bearer ${jwtToken}` },
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      setList(Array.isArray(data) ? data : []);
    } catch (err: any) {
      setError(err.message || 'Failed to load policyholders from Go backend');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPolicyholders();
  }, []);

  return (
    <div className="tab-pane active">
      <div className="card">
        <div className="card-header">
          <h2>Policyholders & Active Policies (Go Backend)</h2>
          <button className="btn" onClick={onOpenCreate}>
            + Create Policyholder
          </button>
        </div>

        <table className="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Full Name</th>
              <th>Email</th>
              <th>Phone</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={6} style={{ textAlign: 'center', color: 'var(--text-muted)' }}>
                  Loading Policyholders from Go Service...
                </td>
              </tr>
            ) : error ? (
              <tr>
                <td colSpan={6} style={{ textAlign: 'center', color: '#F87171' }}>
                  {error}
                </td>
              </tr>
            ) : list.length === 0 ? (
              <tr>
                <td colSpan={6} style={{ textAlign: 'center', color: 'var(--text-muted)' }}>
                  No policyholders found. Click "+ Create Policyholder" to add one!
                </td>
              </tr>
            ) : (
              list.map((p) => (
                <tr key={p.id}>
                  <td style={{ fontFamily: 'JetBrains Mono', fontSize: '11px' }}>
                    {p.id ? p.id.substring(0, 8) + '...' : 'N/A'}
                  </td>
                  <td>
                    <strong>{p.full_name || p.fullName || 'N/A'}</strong>
                  </td>
                  <td>{p.email || 'N/A'}</td>
                  <td>{p.phone || 'N/A'}</td>
                  <td>
                    <span className="badge badge-active">Active</span>
                  </td>
                  <td>
                    <button
                      className="btn btn-secondary"
                      style={{ padding: '4px 10px', fontSize: '11px' }}
                      onClick={() => alert(`Policyholder details:\n${JSON.stringify(p, null, 2)}`)}
                    >
                      Inspect
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};
