import React, { useState } from 'react';
import type { RoutingMode, InspectionResult } from '../types';

const jwtToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMyJ9.ah_KrNhrpocJD5ZKfyjoKMdCZMst4lvLSayT5B_d8Bk";

export const StranglerVisualizer: React.FC = () => {
  const [routingMode, setRoutingMode] = useState<RoutingMode>('auto');
  const [inspection, setInspection] = useState<InspectionResult | null>(null);

  const executeApiTest = async (endpoint: string) => {
    const startTime = performance.now();
    try {
      const headers: Record<string, string> = { Authorization: `Bearer ${jwtToken}` };
      if (routingMode !== 'auto') {
        headers['X-Force-Backend'] = routingMode;
      }

      const res = await fetch(endpoint, { headers });
      const duration = Math.round(performance.now() - startTime);
      const routedTo = res.headers.get('X-Routed-To') || (endpoint.includes('legacy') ? 'legacy-nest (:3001)' : 'go-service (:8080)');
      const targetURL = res.headers.get('X-Target-URL') || (routedTo.includes('3001') ? 'http://localhost:3001' : 'http://localhost:8080');

      let bodyText = '';
      const contentType = res.headers.get('Content-Type') || '';
      if (contentType.includes('application/json')) {
        try {
          const data = await res.json();
          bodyText = JSON.stringify(data, null, 2);
        } catch {
          bodyText = await res.text();
        }
      } else {
        bodyText = await res.text();
      }

      setInspection({
        endpoint,
        status: res.ok ? `${res.status} OK` : `${res.status} ${res.statusText}`,
        ok: res.ok,
        routedTo,
        targetURL,
        latencyMs: duration,
        body: bodyText,
      });
    } catch (err: any) {
      setInspection({
        endpoint,
        status: 'Network Error',
        ok: false,
        routedTo: 'N/A',
        targetURL: 'N/A',
        latencyMs: 0,
        body: err.message || 'Failed to connect to gateway',
      });
    }
  };

  const handleModeChange = (mode: RoutingMode) => {
    setRoutingMode(mode);
  };

  return (
    <div className="tab-pane active">
      <div className="card">
        <div className="card-header flex-wrap">
          <h2>🔀 Live Strangler-Fig Gateway Traffic Map</h2>
          <div className="button-group">
            <button
              className={`tab-btn ${routingMode === 'auto' ? 'active' : ''}`}
              onClick={() => handleModeChange('auto')}
            >
              🔀 Auto (Strangler Rules)
            </button>
            <button
              className={`tab-btn green-btn ${routingMode === 'go' ? 'active' : ''}`}
              onClick={() => handleModeChange('go')}
            >
              ⚡ Force Go (:8080)
            </button>
            <button
              className={`tab-btn amber-btn ${routingMode === 'legacy' ? 'active' : ''}`}
              onClick={() => handleModeChange('legacy')}
            >
              🐢 Force Legacy (:3001)
            </button>
          </div>
        </div>

        <div className="architecture-diagram">
          <div className="arch-node gateway">
            <div className="arch-title">Gateway Proxy</div>
            <div className="arch-subtitle">Port :8000 (Strangler Router)</div>
          </div>

          <div
            className={`flow-line ${
              inspection?.routedTo.includes('go-service') || routingMode === 'go'
                ? 'active-go'
                : 'active-nest'
            }`}
          ></div>

          <div className="arch-node go">
            <div className="arch-title">Go Microservice</div>
            <div className="arch-subtitle">Port :8080 (pgx + sqlc)</div>
          </div>

          <div className="arch-node nest">
            <div className="arch-title">Legacy NestJS</div>
            <div className="arch-subtitle">Port :3001 (TypeORM)</div>
          </div>
        </div>
      </div>

      <div className="grid-2">
        <div className="card">
          <div className="card-header">
            <h2>Configured Service Endpoints</h2>
            <button className="btn btn-secondary" onClick={() => executeApiTest('/api/v1/policyholders')}>
              Run Suite
            </button>
          </div>

          <div className="endpoint-list">
            <div className="endpoint-item">
              <div className="endpoint-info">
                <span className="method-tag get">GET</span>
                <span className="endpoint-path">/api/v1/policyholders</span>
              </div>
              <div className="route-target-pill go">X-Routed-To: go-service</div>
              <button className="btn btn-secondary" onClick={() => executeApiTest('/api/v1/policyholders')}>
                Test
              </button>
            </div>

            <div className="endpoint-item">
              <div className="endpoint-info">
                <span className="method-tag get">GET</span>
                <span className="endpoint-path">/api/v1/policies/expiring?within_days=30</span>
              </div>
              <div className="route-target-pill go">X-Routed-To: go-service</div>
              <button
                className="btn btn-secondary"
                onClick={() => executeApiTest('/api/v1/policies/expiring?within_days=30')}
              >
                Test
              </button>
            </div>

            <div className="endpoint-item">
              <div className="endpoint-info">
                <span className="method-tag post">POST</span>
                <span className="endpoint-path">/api/v1/policyholders</span>
              </div>
              <div className="route-target-pill go">X-Routed-To: go-service</div>
              <button className="btn btn-secondary" onClick={() => executeApiTest('/api/v1/policyholders')}>
                Test
              </button>
            </div>

            <div className="endpoint-item">
              <div className="endpoint-info">
                <span className="method-tag get">GET</span>
                <span className="endpoint-path">/api/v1/legacy-contacts</span>
              </div>
              <div className="route-target-pill nest">X-Routed-To: legacy-nest</div>
              <button className="btn btn-secondary" onClick={() => executeApiTest('/api/v1/legacy-contacts')}>
                Test
              </button>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="card-header">
            <h2>Live Gateway Inspector Output</h2>
            <span className="badge badge-active">
              Latency: {inspection ? `${inspection.latencyMs} ms` : '-- ms'}
            </span>
          </div>

          <div className="console-box">
            {!inspection ? (
              <>
                <div className="console-header-line">// Click "Test" on any endpoint above to trigger a live Gateway request</div>
                <div>[System] Gateway active on http://localhost:8000</div>
              </>
            ) : (
              <>
                <div className="console-header-line">
                  HTTP GET {inspection.endpoint}{' '}
                  {routingMode !== 'auto' && <span style={{ color: '#A78BFA' }}>[forced → {routingMode}]</span>}
                </div>
                <div>
                  <span className="console-key">Status:</span>{' '}
                  <span style={{ color: inspection.ok ? '#34D399' : '#F87171' }}>{inspection.status}</span>
                </div>
                <div>
                  <span className="console-key">X-Routed-To:</span>{' '}
                  <span className={inspection.routedTo.includes('go-service') ? 'console-val-go' : 'console-val-nest'}>
                    {inspection.routedTo}
                  </span>
                </div>
                <div>
                  <span className="console-key">Target Destination:</span>{' '}
                  <span style={{ color: '#FBBF24', fontFamily: 'JetBrains Mono' }}>{inspection.targetURL}</span>
                </div>
                <div>
                  <span className="console-key">Response Body:</span>
                </div>
                <pre style={{ color: inspection.ok ? '#A7F3D0' : '#F87171', fontSize: '12px' }}>{inspection.body}</pre>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
