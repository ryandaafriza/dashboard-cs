import React from 'react';
import { formatDelta } from '../utils/formatUtils';

export function SLACard({ slaToday, slaDelta, slaStatus }) {
  const isNegative = slaDelta < 0;
  const hasDelta   = slaDelta !== 0 && slaDelta !== null && slaDelta !== undefined;

  return (
    <div className="card sla-card animate-in">
      <div className="card__inner">
        <div className="card-header">
          <span className="card-label">SLA Today</span>
          <span className="card-icon card-icon--success">✓</span>
        </div>

        {/* Center content vertically */}
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div className="sla-value">{slaToday.toFixed(2)}%</div>

          {hasDelta && (
            <div className={`sla-delta ${isNegative ? 'sla-delta--negative' : 'sla-delta--positive'}`}>
              <span>{isNegative ? '▼' : '▲'}</span>
              <span>{formatDelta(slaDelta)}</span>
            </div>
          )}

          <div className="sla-status">{slaStatus}</div>
        </div>

        <div className="sla-progress-bar">
          <div className="progress-track">
            <div
              className="progress-fill"
              style={{ width: `${Math.min(slaToday, 100)}%` }}
            />
          </div>
        </div>
      </div>
    </div>
  );
}