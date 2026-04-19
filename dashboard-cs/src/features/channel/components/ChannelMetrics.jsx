import React from 'react';
import { getSLAStatus } from '../utils/channelUtils';

export function ChannelMetrics({ sla, open, closed }) {
  const status = getSLAStatus(sla);

  return (
    <div className="channel-metrics">
      {/* SLA */}
      <div className="channel-metric">
        <span className="channel-metric__label">SLA</span>
        <span
          className="channel-metric__value channel-metric__value--sla"
          style={{
            color: status.color,
            background: status.bg,
            border: `1px solid ${status.border}`,
          }}
        >
          {sla.toFixed(sla % 1 === 0 ? 0 : 2)}%
        </span>
      </div>

      {/* Open */}
      <div className="channel-metric">
        <span className="channel-metric__label">Open</span>
        <span
          className="channel-metric__value"
          style={{
            color: open > 0 ? '#ef4444' : 'var(--text-muted)',
            background: open > 0 ? 'rgba(239,68,68,0.08)' : 'rgba(255,255,255,0.03)',
            border: `1px solid ${open > 0 ? 'rgba(239,68,68,0.2)' : 'rgba(255,255,255,0.06)'}`,
          }}
        >
          {open}
        </span>
      </div>

      {/* Closed */}
      <div className="channel-metric">
        <span className="channel-metric__label">Closed</span>
        <span
          className="channel-metric__value"
          style={{
            color: '#22c55e',
            background: 'rgba(34,197,94,0.08)',
            border: '1px solid rgba(34,197,94,0.2)',
          }}
        >
          {closed}
        </span>
      </div>
    </div>
  );
}
