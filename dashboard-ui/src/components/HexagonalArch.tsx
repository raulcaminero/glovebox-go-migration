import React from 'react';
import { ExternalLink } from 'lucide-react';
import type { AdrInfo } from '../types';

interface HexagonalArchProps {
  onOpenDoc: (docId: string, title: string) => void;
}

const adrs: AdrInfo[] = [
  {
    id: 'adr-0001',
    title: 'ADR 0001: Rewrite Backend in Go',
    summary: 'Why migrate from NestJS to Go incrementally using strangler-fig pattern.',
    status: 'Accepted',
  },
  {
    id: 'adr-0002',
    title: 'ADR 0002: pgx + sqlc over ORM',
    summary: 'Why raw SQL + compile-time generated Go code beat GORM for search & queries.',
    status: 'Accepted',
  },
  {
    id: 'adr-0003',
    title: 'ADR 0003: Strangler-Fig Gateway',
    summary: 'Why a reverse proxy routing layer beats a high-risk big-bang cutover.',
    status: 'Accepted',
  },
  {
    id: 'adr-0004',
    title: 'ADR 0004: Shared Postgres DB Coexistence',
    summary: 'Shared Postgres data layer during strangling vs complex CDC / Debezium dual-writes.',
    status: 'Accepted',
  },
  {
    id: 'adr-0005',
    title: 'ADR 0005: Shadow Traffic & Canaries',
    summary: 'Dark launching read requests and OpenTelemetry metric verification prior to live cutover.',
    status: 'Accepted',
  },
  {
    id: 'adr-0006',
    title: 'ADR 0006: React 18 + TypeScript Dashboard UI',
    summary: 'Modular React + Vite component architecture embedded into Go gateway via //go:embed static/*.',
    status: 'Accepted',
  },
];

export const HexagonalArch: React.FC<HexagonalArchProps> = ({ onOpenDoc }) => {
  return (
    <div className="tab-pane active">
      {/* Hexagonal Layer Diagram Card */}
      <div className="card">
        <div className="card-header">
          <h2>Hexagonal Architecture (Ports & Adapters) — Layered Design</h2>
          <div className="flex-row gap-2">
            <a
              href="https://github.com/raulcaminero/glovebox-go-migration"
              target="_blank"
              rel="noreferrer"
              className="btn btn-secondary text-xs"
            >
              🐙 Open Repository on GitHub <ExternalLink size={12} />
            </a>
            <span className="badge badge-active">Strict Boundary Isolation</span>
          </div>
        </div>

        <div className="flex-col gap-4">
          <p className="text-muted text-sm">
            The Go service follows a strict <strong>Hexagonal (Clean) Architecture</strong> pattern. Core business domain logic is completely isolated from database drivers, HTTP frameworks, and third-party libraries.
          </p>

          <div className="hex-grid">
            {/* Layer 1: Adapters */}
            <div className="hex-layer layer-cyan">
              <div className="layer-tag text-cyan">1. ADAPTERS (Driving / Driven)</div>
              <h4>Transport & Persistence</h4>
              <ul>
                <li><code>transport/http</code> (Chi HTTP Handlers)</li>
                <li><code>repo/policy_pg.go</code> (Postgres Adapter)</li>
                <li><code>repo/sqlcgen</code> (Generated SQL Code)</li>
              </ul>
              <div className="layer-footer">Implements ports & handles I/O</div>
            </div>

            {/* Layer 2: Ports & Services */}
            <div className="hex-layer layer-purple">
              <div className="layer-tag text-purple">2. PORTS & SERVICES</div>
              <h4>Use Cases & Interfaces</h4>
              <ul>
                <li><code>service/policy.go</code> (Business Rules)</li>
                <li><code>domain.PolicyRepo</code> (Repository Port)</li>
                <li><code>domain.PolicyholderRepo</code> (Interface)</li>
              </ul>
              <div className="layer-footer">Orchestrates business logic & ports</div>
            </div>

            {/* Layer 3: Core Domain */}
            <div className="hex-layer layer-green">
              <div className="layer-tag text-green">3. CORE DOMAIN (Zero Dependencies)</div>
              <h4>Pure Domain Entities</h4>
              <ul>
                <li><code>domain/policy.go</code> (Policy Entity)</li>
                <li><code>domain/policyholder.go</code> (Entity)</li>
                <li>Zero external imports (no pgx, sqlc, chi)</li>
              </ul>
              <div className="layer-footer">Pure Go types & validation</div>
            </div>
          </div>

          <div className="data-flow-bar font-mono text-xs">
            <span className="text-cyan">HTTP Request</span> ➔{' '}
            <span className="text-purple">Chi Handler</span> ➔{' '}
            <span style={{ color: '#38BDF8' }}>PolicyService</span> ➔{' '}
            <span style={{ color: '#FBBF24' }}>PolicyRepo Interface</span> ➔{' '}
            <span style={{ color: '#10B981' }}>SQLC Postgres Adapter</span> ➔{' '}
            <span style={{ color: '#34D399' }}>PostgreSQL</span>
          </div>
        </div>
      </div>

      <div className="section-title">
        <span>Architecture Decision Records (ADRs)</span>
      </div>

      <div className="grid-3">
        {adrs.map((adr) => (
          <div
            key={adr.id}
            className="adr-card clickable"
            onClick={() => onOpenDoc(adr.id, adr.title)}
          >
            <div className="adr-tag">{adr.id.toUpperCase()} 📄 (Click to Read Full Doc)</div>
            <h3>{adr.title.split(': ')[1] || adr.title}</h3>
            <p className="text-muted text-xs">{adr.summary}</p>
            <div className="status-label text-cyan">Status: {adr.status}</div>
          </div>
        ))}
      </div>
    </div>
  );
};
