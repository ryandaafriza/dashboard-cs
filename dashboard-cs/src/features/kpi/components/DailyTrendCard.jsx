import React, { useMemo, useState } from 'react';

const W       = 320;
const H       = 80;
const PAD     = { top: 8, right: 8, bottom: 4, left: 8 };
const CHART_W = W - PAD.left - PAD.right;
const CHART_H = H - PAD.top - PAD.bottom;

function formatDate(dateStr) {
  // "2026-04-13T00:00:00+07:00" → "13/04"
  const d = new Date(dateStr);
  const day   = String(d.getDate()).padStart(2, '0');
  const month = String(d.getMonth() + 1).padStart(2, '0');
  return `${day}/${month}`;
}

function formatDateFull(dateStr) {
  const d = new Date(dateStr);
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
}

export function DailyTrendCard({ data }) {
  const [tooltip, setTooltip] = useState(null); // { x, y, item }

  const { points, linePath, areaPath, yMax } = useMemo(() => {
    if (!data || data.length < 1) return { points: [], linePath: '', areaPath: '', yMax: 0 };

    const values = data.map((d) => d.created ?? d.value ?? 0);
    const max    = Math.max(...values, 1);

    const pts = data.map((d, i) => {
      const val = d.created ?? d.value ?? 0;
      const x   = PAD.left + (data.length === 1 ? CHART_W / 2 : (i / (data.length - 1)) * CHART_W);
      const y   = PAD.top  + (1 - val / max) * CHART_H;
      return { x, y, val, date: d.date ?? d.time };
    });

    // Smooth bezier line
    const line = pts.reduce((acc, p, i) => {
      if (i === 0) return `M ${p.x} ${p.y}`;
      const prev = pts[i - 1];
      const cpx  = (prev.x + p.x) / 2;
      return `${acc} C ${cpx} ${prev.y}, ${cpx} ${p.y}, ${p.x} ${p.y}`;
    }, '');

    const area = pts.length < 2 ? '' :
      `${line} L ${pts[pts.length - 1].x} ${H} L ${pts[0].x} ${H} Z`;

    return { points: pts, linePath: line, areaPath: area, yMax: max };
  }, [data]);

  // Tampilkan max 4 label tanggal di sumbu X agar tidak penuh
  const xLabels = useMemo(() => {
    if (!points.length) return [];
    if (points.length <= 4) return points;
    const step = Math.floor((points.length - 1) / 3);
    return [0, step, step * 2, points.length - 1].map((i) => points[i]);
  }, [points]);

  function handleMouseMove(e) {
    const svgEl  = e.currentTarget;
    const rect   = svgEl.getBoundingClientRect();
    const mouseX = ((e.clientX - rect.left) / rect.width) * W;

    // Cari titik terdekat
    let closest = null, minDist = Infinity;
    points.forEach((p) => {
      const dist = Math.abs(p.x - mouseX);
      if (dist < minDist) { minDist = dist; closest = p; }
    });

    if (closest && minDist < 30) {
      setTooltip(closest);
    } else {
      setTooltip(null);
    }
  }

  return (
    <div className="card chart-card animate-in">
      <div className="card__inner">
        <div className="card-header">
          <span className="card-label">Daily Trend</span>
          {tooltip && (
            <span style={{ fontSize: 11, color: 'var(--text-muted)', fontFamily: 'var(--font-mono)' }}>
              {formatDateFull(tooltip.date)}
            </span>
          )}
        </div>

        <div className="chart-area" style={{ position: 'relative' }}>
          <svg
            viewBox={`0 0 ${W} ${H}`}
            preserveAspectRatio="none"
            style={{ width: '100%', height: '100%', overflow: 'visible' }}
            onMouseMove={handleMouseMove}
            onMouseLeave={() => setTooltip(null)}
          >
            <defs>
              <linearGradient id="trendGrad2" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%"   stopColor="#3b82f6" stopOpacity="0.2" />
                <stop offset="100%" stopColor="#3b82f6" stopOpacity="0" />
              </linearGradient>
            </defs>

            {/* Area fill */}
            {areaPath && <path d={areaPath} fill="url(#trendGrad2)" />}

            {/* Line */}
            {linePath && (
              <path
                d={linePath} fill="none"
                stroke="#3b82f6" strokeWidth="1.8"
                strokeLinecap="round" strokeLinejoin="round"
              />
            )}

            {/* Dots pada setiap titik data */}
            {points.map((p, i) => (
              <circle
                key={i} cx={p.x} cy={p.y} r="2.5"
                fill={tooltip?.date === p.date ? '#3b82f6' : 'var(--bg-card)'}
                stroke="#3b82f6" strokeWidth="1.5"
              />
            ))}

            {/* Tooltip crosshair + bubble */}
            {tooltip && (
              <>
                {/* Vertical line */}
                <line
                  x1={tooltip.x} y1={PAD.top}
                  x2={tooltip.x} y2={H - PAD.bottom}
                  stroke="#3b82f6" strokeWidth="1"
                  strokeDasharray="3 2" strokeOpacity="0.5"
                />

                {/* Tooltip bubble — posisi dinamis agar tidak keluar kiri/kanan */}
                {(() => {
                  const bw = 72, bh = 32, br = 5;
                  const bx = Math.min(Math.max(tooltip.x - bw / 2, PAD.left), W - PAD.right - bw);
                  const by = PAD.top - bh - 4;
                  return (
                    <g>
                      <rect
                        x={bx} y={by} width={bw} height={bh}
                        rx={br} ry={br}
                        fill="var(--bg-elevated, #242d42)"
                        stroke="rgba(59,130,246,0.3)" strokeWidth="1"
                      />
                      <text
                        x={bx + bw / 2} y={by + 11}
                        textAnchor="middle" fontSize="9"
                        fill="var(--text-muted, #8a95a8)"
                        fontFamily="var(--font-mono)"
                      >
                        {formatDate(tooltip.date)}
                      </text>
                      <text
                        x={bx + bw / 2} y={by + 24}
                        textAnchor="middle" fontSize="12"
                        fontWeight="600" fill="#3b82f6"
                        fontFamily="var(--font-mono)"
                      >
                        {tooltip.val.toLocaleString('id-ID')}
                      </text>
                    </g>
                  );
                })()}
              </>
            )}
          </svg>
        </div>

        {/* Sumbu X — label tanggal */}
        <div style={{
          display: 'flex', justifyContent: 'space-between',
          marginTop: 4, paddingInline: PAD.left,
        }}>
          {xLabels.map((p, i) => (
            <span key={i} style={{
              fontSize: 10, color: 'var(--text-muted)',
              fontFamily: 'var(--font-mono)',
            }}>
              {formatDate(p.date)}
            </span>
          ))}
        </div>
      </div>
    </div>
  );
}