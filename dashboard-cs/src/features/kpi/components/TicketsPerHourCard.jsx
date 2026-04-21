import React, { useMemo, useState } from 'react';
import { normalizeBarData } from '../utils/chartUtils';

export function TicketsPerHourCard({ data }) {
  const [activeIdx, setActiveIdx] = useState(null);

  const createdNorm = useMemo(() => normalizeBarData(data, 'created'), [data]);
  const solvedNorm  = useMemo(() => normalizeBarData(data, 'solved'),  [data]);

  const BAR_H   = 56;
  const BAR_W   = 10;
  const GAP     = 3;
  const GROUP_W = BAR_W * 2 + GAP + 4;
  const TOTAL_W = data.length * GROUP_W;

  function handleMouseMove(e) {
    const svgEl  = e.currentTarget;
    const rect   = svgEl.getBoundingClientRect();
    const mouseX = ((e.clientX - rect.left) / rect.width) * TOTAL_W;

    let closestIdx = -1;
    let minDist = Infinity;
    data.forEach((_, i) => {
      const dist = Math.abs(i * GROUP_W + GROUP_W / 2 - mouseX);
      if (dist < minDist) { minDist = dist; closestIdx = i; }
    });

    setActiveIdx(closestIdx >= 0 && minDist < GROUP_W ? closestIdx : null);
  }

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

        <div className="chart-area" style={{ overflowX: 'auto', position: 'relative' }}>
          <svg
            viewBox={`0 0 ${TOTAL_W} ${BAR_H + 16}`}
            preserveAspectRatio="none"
            style={{ width: '100%', height: '100%', minHeight: 72 }}
            onMouseMove={handleMouseMove}
            onMouseLeave={() => setActiveIdx(null)}
          >
            {createdNorm.map((d, i) => {
              const x        = i * GROUP_W;
              const createdH = Math.max(2, d.normalized * BAR_H);
              const solvedH  = Math.max(2, solvedNorm[i].normalized * BAR_H);
              const isActive = activeIdx === i;

              const createdVal = data[i]?.created ?? 0;
              const solvedVal  = data[i]?.solved  ?? 0;
              const MIN_H_FOR_INSIDE = 14;

              return (
                <g key={d.hour}>
                  {isActive && (
                    <rect
                      x={x - 2} y={0}
                      width={GROUP_W} height={BAR_H}
                      rx={3}
                      fill="rgba(255,255,255,0.06)"
                    />
                  )}

                  <rect
                    x={x} y={BAR_H - createdH}
                    width={BAR_W} height={createdH}
                    rx={2}
                    fill={isActive ? 'rgba(59,130,246,1)' : 'rgba(59,130,246,0.7)'}
                  />

                  <rect
                    x={x + BAR_W + GAP} y={BAR_H - solvedH}
                    width={BAR_W} height={solvedH}
                    rx={2}
                    fill={isActive ? 'rgba(245,158,11,1)' : 'rgba(245,158,11,0.7)'}
                  />

                  {isActive && createdVal > 0 && (
                    createdH >= MIN_H_FOR_INSIDE ? (
                      <text
                        x={x + BAR_W / 2} y={BAR_H - createdH / 2}
                        fontSize={7} fill="white"
                        textAnchor="middle" dominantBaseline="middle"
                        fontFamily="var(--font-mono)" fontWeight="700"
                        transform={`rotate(-90, ${x + BAR_W / 2}, ${BAR_H - createdH / 2})`}
                      >
                        {createdVal}
                      </text>
                    ) : (
                      <text
                        x={x + BAR_W / 2} y={BAR_H - createdH - 3}
                        fontSize={7} fill="#3b82f6"
                        textAnchor="middle"
                        fontFamily="var(--font-mono)" fontWeight="700"
                      >
                        {createdVal}
                      </text>
                    )
                  )}

                  {isActive && solvedVal > 0 && (
                    solvedH >= MIN_H_FOR_INSIDE ? (
                      <text
                        x={x + BAR_W + GAP + BAR_W / 2} y={BAR_H - solvedH / 2}
                        fontSize={7} fill="white"
                        textAnchor="middle" dominantBaseline="middle"
                        fontFamily="var(--font-mono)" fontWeight="700"
                        transform={`rotate(-90, ${x + BAR_W + GAP + BAR_W / 2}, ${BAR_H - solvedH / 2})`}
                      >
                        {solvedVal}
                      </text>
                    ) : (
                      <text
                        x={x + BAR_W + GAP + BAR_W / 2} y={BAR_H - solvedH - 3}
                        fontSize={7} fill="#f59e0b"
                        textAnchor="middle"
                        fontFamily="var(--font-mono)" fontWeight="700"
                      >
                        {solvedVal}
                      </text>
                    )
                  )}

                  {i % 2 === 0 && (
                    <text
                      x={x + BAR_W} y={BAR_H + 12}
                      fontSize={7}
                      fill={isActive ? 'var(--text-primary, #e2e8f0)' : 'var(--text-muted)'}
                      textAnchor="middle" fontFamily="var(--font-mono)"
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