# ADR 0006: React 18 + TypeScript + Vite architecture for Migration Control Center Dashboard UI

## Status
Accepted

## Context
The Job Description requires **"Modern Frontend Proficiency (React, Vue)"** alongside Go backend engineering skills. While the migration control gateway acts as a backend proxy, operating a production strangler-fig rollout requires real-time routing inspection, live CRUD testing, and metric monitoring. 

Relying solely on a single static HTML file leaves the frontend story under-documented and misses an opportunity to demonstrate production React architecture.

## Decision
Build the **Migration Control Center Dashboard UI** as a modular **React 18 + TypeScript + Vite** application residing in `dashboard-ui/`.

The build pipeline compiles the React app into `gateway/static/`, which Go embeds directly via `//go:embed static/*` and serves via `http.FileServer`.

Key Component Architecture:
- `src/components/Header.tsx`: Real-time backend status indicators and repository links.
- `src/components/Navigation.tsx`: Tab-based navigation control.
- `src/components/StranglerVisualizer.tsx`: Interactive gateway proxy tester, header overrides, and live traffic inspector.
- `src/components/HexagonalArch.tsx`: Visual breakdown of Ports & Adapters layers and request pipelines.
- `src/components/CrmManager.tsx`: Real-time Policyholder CRUD table interacting with Go REST APIs.
- `src/components/AdrViewer.tsx`: Architecture Decision Record reader and modal renderer.
- `src/components/AiGuardrails.tsx`: CLAUDE.md rules and `tools/ai-pr-review` guardrail breakdown.

## Alternatives Considered

### 1. Embedded Single Vanilla HTML/JS File
- **Approach**: Maintain a single static HTML file with inline vanilla JavaScript.
- **Why Rejected**: Does not demonstrate modern component architecture, TypeScript type safety, or state management; fails to address the JD's explicit React/Vue requirement.

### 2. Vue.js 3 + Vite Application
- **Approach**: Build the dashboard in Vue 3 with Composition API.
- **Why Deferred**: React 18 was chosen due to stronger team alignment, but Vue 3 remains a compatible alternative for future micro-frontend modules.

## Consequences
+ **Fulfills JD Frontend Requirement**: Proves React + TypeScript proficiency directly within the migration repository.
+ **Clean Component Separation**: Decoupled, reusable UI components with strict TypeScript types (`src/types.ts`).
+ **Zero Operational Overhead**: The compiled React bundle is embedded directly inside the Go Gateway binary (`gateway`), preserving single-binary deployment simplicity (`go run .`).
