import React from 'react';
import { formatNumber } from '../../kpi/utils/formatUtils';

export function SummarySection({ summary }) {
  const { total, open, closed } = summary;
  const openPct   = total > 0 ? ((open   / total) * 100).toFixed(1) : '0.0';
  const closedPct = total > 0 ? ((closed / total) * 100).toFixed(1) : '0.0';

  return (
    <section>
      <div className="section-label">Volume Summary</div>
      <div className="summary-section">
        {/* Left: totals */}
        <div className="card summary-card animate-in">
          <div className="card__inner" style={{ flexDirection: 'column', gap: 'var(--space-3)', padding: 'var(--space-4) var(--space-5)' }}>
            <div className="card-header" style={{ marginBottom: 0 }}>
              <span className="card-label">Total Tickets</span>
            </div>
            <div className="summary-total-value">{formatNumber(total)}</div>

            <div style={{ display: 'flex', gap: 'var(--space-5)', marginTop: 'var(--space-1)' }}>
              <div className="summary-stat">
                <span
                  className="summary-stat-value summary-stat-value--open"
                >
                  {formatNumber(open)}
                </span>
                <span className="summary-stat-label">Open</span>
              </div>

              <div
                style={{
                  width: 1,
                  background: 'var(--border-default)',
                  alignSelf: 'stretch',
                }}
              />

              <div className="summary-stat">
                <span className="summary-stat-value summary-stat-value--closed">
                  {formatNumber(closed)}
                </span>
                <span className="summary-stat-label">Closed</span>
              </div>
            </div>
          </div>
        </div>

        {/* Right: bar breakdown */}
        <div className="card animate-in" style={{ animationDelay: '60ms' }}>
          <div
            className="card__inner"
            style={{ justifyContent: 'center', gap: 'var(--space-4)' }}
          >
            <div className="card-header">
              <span className="card-label">Ticket Status Breakdown</span>
              <span style={{ fontSize: 11, color: 'var(--text-muted)', fontFamily: 'var(--font-mono)' }}>
                Total: {formatNumber(total)}
              </span>
            </div>

            <div className="summary-bar-section" style={{ flex: 1 }}>
              {/* Open bar */}
              <div className="summary-bar-row">
                <div className="summary-bar-label-row">
                  <span className="summary-bar-label">
                    <span
                      style={{
                        display: 'inline-block',
                        width: 8,
                        height: 8,
                        borderRadius: 2,
                        background: 'var(--color-danger)',
                        marginRight: 6,
                        verticalAlign: 'middle',
                      }}
                    />
                    Open
                  </span>
                  <span className="summary-bar-pct">
                    {formatNumber(open)} ({openPct}%)
                  </span>
                </div>
                <div className="summary-bar-track" style={{ overflow: 'visible' }}>
                  <div
                    className="summary-bar-fill summary-bar-fill--open"
                    style={{ width: parseFloat(openPct) === 0 ? '4px' : `${openPct}%`, opacity: parseFloat(openPct) === 0 ? 0.35 : 1 }}
                  />
                </div>
              </div>

              {/* Closed bar */}
              <div className="summary-bar-row">
                <div className="summary-bar-label-row">
                  <span className="summary-bar-label">
                    <span
                      style={{
                        display: 'inline-block',
                        width: 8,
                        height: 8,
                        borderRadius: 2,
                        background: 'var(--color-success)',
                        marginRight: 6,
                        verticalAlign: 'middle',
                      }}
                    />
                    Closed
                  </span>
                  <span className="summary-bar-pct">
                    {formatNumber(closed)} ({closedPct}%)
                  </span>
                </div>
                <div className="summary-bar-track" style={{ overflow: 'visible' }}>
                  <div
                    className="summary-bar-fill summary-bar-fill--closed"
                    style={{ width: parseFloat(closedPct) === 0 ? '4px' : `${closedPct}%`, opacity: parseFloat(closedPct) === 0 ? 0.35 : 1 }}
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}