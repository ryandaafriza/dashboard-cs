import React from 'react';
import { PriorityBadge } from './PriorityBadge';

export function PrioritySection({ priority }) {
  const items = [
    { type: 'roaming',     count: priority.roaming },
    { type: 'extra_quota', count: priority.extra_quota },
    { type: 'cc',          count: priority.cc },
    { type: 'vip',         count: priority.vip },
    { type: 'p1',          count: priority.p1 },
    { type: 'urgent',      count: priority.urgent },
  ];

  const totalActive = items.reduce((s, i) => s + i.count, 0);
  const criticalCount = priority.p1 + priority.urgent;

  return (
    <section>
      <div className="section-label">Priority Tickets</div>
      <div className="card priority-section-card animate-in">
        <div className="card__inner priority-section-inner">

          {/* Left: header info */}
          <div className="priority-section-meta">
            <div className="priority-meta-title">
              Unresolved Priority
            </div>
            <div className="priority-meta-total">
              {totalActive}
              <span className="priority-meta-sub">active</span>
            </div>
            {criticalCount > 0 && (
              <div className="priority-critical-alert">
                <span className="priority-critical-dot" />
                {criticalCount} critical needs immediate action
              </div>
            )}
          </div>

          <div className="priority-divider" />

          {/* Right: badges grid */}
          <div className="priority-badges-grid">
            {items.map((item) => (
              <PriorityBadge key={item.type} type={item.type} count={item.count} />
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
