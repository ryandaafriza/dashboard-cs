/**
 * Format number with dot separator (Indonesian style)
 * e.g. 1295 → "1.295"
 */
export function formatNumber(n) {
  return n?.toLocaleString('id-ID') ?? '—';
}

/**
 * Format percentage with sign
 */
export function formatDelta(delta) {
  const sign = delta >= 0 ? '+' : '';
  return `${sign}${delta.toFixed(2)}%`;
}

/**
 * Format time HH:MM
 */
export function formatTime(date) {
  return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
}
