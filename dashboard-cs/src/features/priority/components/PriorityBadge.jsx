import React from 'react';
import { PRIORITY_CONFIG, getUrgencyLevel } from '../utils/priorityUtils';

export function PriorityBadge({ type, count }) {
  const config = PRIORITY_CONFIG[type] ?? {
    label: type,
    color: '#8a95a8',
    bg: 'rgba(138,149,168,0.1)',
    border: 'rgba(138,149,168,0.2)',
    icon: '•',
  };

  const urgency = getUrgencyLevel(type);
  const isEmpty = count === 0;
  const isCritical = urgency === 'critical' && count > 0;

  return (
    <div
      className={`priority-badge ${isCritical ? 'priority-badge--critical' : ''}`}
      style={{
        '--p-color': config.color,
        '--p-bg': isEmpty ? 'rgba(255,255,255,0.03)' : config.bg,
        '--p-border': isEmpty ? 'rgba(255,255,255,0.06)' : config.border,
      }}
    >
      <div className="priority-badge__header">
        <span className="priority-badge__icon">{config.icon}</span>
        <span className="priority-badge__label">{config.label}</span>
      </div>
      <div className={`priority-badge__count ${isEmpty ? 'priority-badge__count--empty' : ''}`}>
        {count}
      </div>
      {isCritical && <div className="priority-badge__pulse" />}
    </div>
  );
}
