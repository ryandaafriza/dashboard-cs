import React from 'react';
import { formatNumber } from '../utils/formatUtils';

const centerStyle = {
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'center',
};

/* ─────────────────────────────────────────
   Created Today Card
───────────────────────────────────────── */
export function CreatedTodayCard({ value, delta }) {
  const hasDelta   = delta !== 0 && delta !== null && delta !== undefined;
  const isNegative = delta < 0;

  return (
    <div className="card stat-card animate-in">
      <div className="card__inner">
        <div className="card-header">
          <span className="card-label">Created Today</span>
          <span className="card-icon card-icon--info">📋</span>
        </div>
        <div style={centerStyle}>
          <div className="stat-value">{formatNumber(value)}</div>
          {hasDelta && (
            <div className={`stat-meta ${isNegative ? 'stat-meta--negative' : 'stat-meta--positive'}`}>
              <span>{isNegative ? '▼' : '▲'}</span>
              <span>{isNegative ? '' : '+'}{formatNumber(delta)} vs last day</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

/* ─────────────────────────────────────────
   Open Tickets Card
───────────────────────────────────────── */
export function OpenTicketsCard({ value }) {
  return (
    <div className="card stat-card animate-in">
      <div className="card__inner">
        <div className="card-header">
          <span className="card-label">Open Tickets</span>
          <span className="card-icon card-icon--warning">⏳</span>
        </div>
        <div style={centerStyle}>
          <div className="stat-value" style={{ color: 'var(--color-warning)' }}>
            {formatNumber(value)}
          </div>
          <div className="stat-meta">Active</div>
        </div>
      </div>
    </div>
  );
}

/* ─────────────────────────────────────────
   Unassigned Card
───────────────────────────────────────── */
export function UnassignedCard({ value }) {
  const hasValue = value !== null && value !== undefined;

  return (
    <div className="card stat-card unassigned-card animate-in">
      <div className="card__inner">
        <div className="card-header">
          <span className="card-label">Unassigned</span>
          <span className="card-icon card-icon--danger unassigned-blink">⚠</span>
        </div>
        <div style={centerStyle}>
          <div className="stat-value" style={{ color: 'var(--color-danger)' }}>
            {hasValue ? value : '—'}
          </div>
          <div className="unassigned-label">Needs attention</div>
        </div>
      </div>
    </div>
  );
}

/* ─────────────────────────────────────────
   CSAT Card
───────────────────────────────────────── */
export function CSATCard({ percentage, score }) {
  const circumference = 2 * Math.PI * 26;
  const offset = circumference - (percentage / 100) * circumference;

  return (
    <div className="card csat-card animate-in">
      <div className="card__inner">
        <div className="card-header">
          <span className="card-label">CSAT</span>
          <span className="card-icon card-icon--success">★</span>
        </div>

        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: 'var(--space-3)',
            flex: 1,
            justifyContent: 'center',
          }}
        >
          <div className="csat-ring-wrap">
            <svg width="72" height="72" viewBox="0 0 60 60">
              <circle cx="30" cy="30" r="26" fill="none" stroke="rgba(255,255,255,0.06)" strokeWidth="5" />
              <circle
                cx="30" cy="30" r="26" fill="none"
                stroke="#22c55e" strokeWidth="5" strokeLinecap="round"
                strokeDasharray={circumference} strokeDashoffset={offset}
                style={{ transition: 'stroke-dashoffset 1.2s ease' }}
              />
            </svg>
            <div className="csat-ring-center">{percentage}%</div>
          </div>

          <div>
            <div className="csat-score">{score.toFixed(2)}</div>
            <div className="csat-sub">/ 5.00 score</div>
          </div>
        </div>
      </div>
    </div>
  );
}