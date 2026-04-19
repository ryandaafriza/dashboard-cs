/**
 * Build an SVG polyline path from an array of {value} points.
 * @param {Array<{value: number}>} data
 * @param {number} width
 * @param {number} height
 * @param {number} [padding=4]
 * @returns {{ path: string, points: Array<{x:number, y:number}> }}
 */
export function buildLinePath(data, width, height, padding = 4) {
  if (!data || data.length < 2) return { path: '', points: [] };

  const values = data.map((d) => d.value);
  const min = Math.min(...values);
  const max = Math.max(...values);
  const range = max - min || 1;

  const points = data.map((d, i) => ({
    x: padding + (i / (data.length - 1)) * (width - padding * 2),
    y: padding + (1 - (d.value - min) / range) * (height - padding * 2),
  }));

  // Smooth bezier
  const path = points.reduce((acc, p, i) => {
    if (i === 0) return `M ${p.x} ${p.y}`;
    const prev = points[i - 1];
    const cpx = (prev.x + p.x) / 2;
    return `${acc} C ${cpx} ${prev.y}, ${cpx} ${p.y}, ${p.x} ${p.y}`;
  }, '');

  return { path, points };
}

/**
 * Build area fill path (path + close to bottom)
 */
export function buildAreaPath(data, width, height, padding = 4) {
  const { path, points } = buildLinePath(data, width, height, padding);
  if (!path) return '';
  const first = points[0];
  const last = points[points.length - 1];
  return `${path} L ${last.x} ${height} L ${first.x} ${height} Z`;
}

/**
 * Normalize bar data for rendering
 */
export function normalizeBarData(data, key) {
  const values = data.map((d) => d[key]);
  const max = Math.max(...values) || 1;
  return data.map((d) => ({ ...d, normalized: d[key] / max }));
}
