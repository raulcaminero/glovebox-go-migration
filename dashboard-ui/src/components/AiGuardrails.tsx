import React from 'react';

interface AiGuardrailsProps {
  onOpenDoc: (docId: string, title: string) => void;
}

export const AiGuardrails: React.FC<AiGuardrailsProps> = ({ onOpenDoc }) => {
  return (
    <div className="tab-pane active">
      <div className="card">
        <div className="card-header">
          <h2>AI Tools, Architecture Guardrails & Automated PR Reviewer</h2>
        </div>

        <div className="grid-2">
          <div
            className="adr-card clickable"
            onClick={() => onOpenDoc('claude', 'CLAUDE.md — AI Guardrails & Architecture Rules')}
          >
            <div className="adr-tag">CLAUDE.MD 📄 (Click to Read Full Rules)</div>
            <h3>CLAUDE.md Rules Enforced</h3>
            <ul className="rule-list">
              <li>Zero framework imports (pgx, sqlc, chi) inside <code>internal/domain</code></li>
              <li>Service layer depends strictly on domain interfaces</li>
              <li>HTTP handlers stay thin (decode request -&gt; call service -&gt; encode)</li>
              <li>Table-driven tests for all new business rules</li>
            </ul>
          </div>

          <div
            className="adr-card clickable"
            onClick={() => onOpenDoc('ai-workflow', 'docs/AI_WORKFLOW.md — AI-Assisted Migration Workflow')}
          >
            <div className="adr-tag">AI_WORKFLOW.MD 📄 (Click to Read Full Doc)</div>
            <h3>Automated PR Audit CLI (<code>tools/ai-pr-review</code>)</h3>
            <p className="text-muted text-xs">
              A custom Go tool that sends PR diffs to the Claude API to automatically check compliance against <code>CLAUDE.md</code> standards before human code review.
            </p>
          </div>
        </div>

        <div className="grid-2 mt-4">
          <div
            className="adr-card clickable"
            onClick={() =>
              onOpenDoc('skill', 'Antigravity Skill — ai-pr-review (.agents/skills/ai-pr-review/SKILL.md)')
            }
          >
            <div className="adr-tag">🤖 ANTIGRAVITY SKILL 📄 (Click to Read Full Skill)</div>
            <h3>Interactive AI PR Review Skill</h3>
            <p className="text-muted text-xs">
              An Antigravity IDE skill that lets any developer say <em>"review my diff"</em> and get instant architecture compliance feedback — no command, no API key required.
            </p>
            <div className="tag-purple mt-2">Scope: This workspace only · Auto-discovered by Antigravity</div>
          </div>

          <div className="adr-card">
            <div className="adr-tag">📊 COVERAGE SUMMARY</div>
            <h3>AI Guardrails at a Glance</h3>
            <ul className="rule-list">
              <li><strong style={{ color: '#34D399' }}>CLAUDE.md</strong> — rules for Claude Code CLI (any engineer)</li>
              <li><strong style={{ color: '#38BDF8' }}>tools/ai-pr-review</strong> — Go CLI, pipeable into CI/GitHub Actions</li>
              <li><strong style={{ color: '#A78BFA' }}>Antigravity Skill</strong> — interactive, no command needed</li>
            </ul>
            <div className="text-muted text-xs mt-3">Rules written once ➔ enforced in 3 places</div>
          </div>
        </div>
      </div>
    </div>
  );
};
