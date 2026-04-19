import React, { useMemo } from 'react';
import { normalizeBarData } from '../utils/chartUtils';

export function TicketsPerHourCard({ data }) {
  const createdNorm = useMemo(() => normalizeBarData(data, 'created'), [data]);
  const solvedNorm = useMemo(() => normalizeBarData(data, 'solved'), [data]);

  const BAR_H = 56;
  const BAR_W = 10;
  const GAP = 3;
  const GROUP_W = BAR_W * 2 + GAP + 4;
  const TOTAL_W = data.length * GROUP_W;

  return (
    <div className="card chart-card animate-in">
      <div className="card__inner">
        <div className="card-header">
          <span className="card-label">Tickets Per Hour</span>
          <div className="chart-legend">
            <div className="legend-item">
              <div className="legend-dot" style={{ background: '#3b82f6' }} />
              <span>Created</span>
            </div>
            <div className="legend-item">
              <div className="legend-dot" style={{ background: '#f59e0b' }} />
              <span>Solved</span>
            </div>
          </div>
        </div>

        <div className="chart-area" style={{ overflowX: 'auto' }}>
          <svg
            viewBox={`0 0 ${TOTAL_W} ${BAR_H + 16}`}
            preserveAspectRatio="none"
            style={{ width: '100%', height: '100%', minHeight: 72 }}
          >
            {createdNorm.map((d, i) => {
              const x = i * GROUP_W;
              const createdH = Math.max(2, d.normalized * BAR_H);
              const solvedH = Math.max(2, solvedNorm[i].normalized * BAR_H);

              return (
                <g key={d.hour}>
                  {/* Created bar */}
                  <rect
                    x={x}
                    y={BAR_H - createdH}
                    width={BAR_W}
                    height={createdH}
                    rx={2}
                    fill="rgba(59,130,246,0.7)"
                  />
                  {/* Solved bar */}
                  <rect
                    x={x + BAR_W + GAP}
                    y={BAR_H - solvedH}
                    width={BAR_W}
                    height={solvedH}
                    rx={2}
                    fill="rgba(245,158,11,0.7)"
                  />
                  {/* Hour label every 2 */}
                  {i % 2 === 0 && (
                    <text
                      x={x + BAR_W}
                      y={BAR_H + 12}
                      fontSize={7}
                      fill="var(--text-muted)"
                      textAnchor="middle"
                      fontFamily="var(--font-mono)"
                    >
                      {d.hour}
                    </text>
                  )}
                </g>
              );
            })}
          </svg>
        </div>
      </div>
    </div>
  );
}
