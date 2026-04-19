import { apiFetch } from './apiClient';

/**
 * Fetch real-time KPI snapshot (no date filter — always "now").
 * @returns {Promise<RealtimeResponse>}
 */
export async function fetchRealtime() {
  return apiFetch('/api/v1/realtime');
}

/**
 * @typedef {Object} RealtimeResponse
 * @property {{ percentage: number, delta: number }} sla_today
 * @property {{ total: number, delta: number }} created_today
 * @property {number} open_tickets
 * @property {number} unassigned
 * @property {number} incidents_active
 */
