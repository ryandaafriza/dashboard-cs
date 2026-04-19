import { apiFetch, buildQuery } from './apiClient';

/**
 * Fetch filtered dashboard data.
 * @param {{ from: string, to: string }} filter  — ISO date strings e.g. "2026-04-13"
 * @returns {Promise<DashboardResponse>}
 */
export async function fetchDashboard({ from, to }) {
  const qs = buildQuery({ from, to });
  return apiFetch(`/api/v1/dashboard${qs}`);
}

/**
 * @typedef {Object} DashboardResponse
 * @property {{ from: string, to: string }} filter
 * @property {{ total_tickets: number, open: number, closed: number, csat_percentage: number, csat_score: number }} summary
 * @property {Array<{ date: string, created: number, solved: number }>} daily_trend
 * @property {Array<{ hour: string, created: number, solved: number }>} tickets_per_hour
 * @property {{ roaming: number, extra_quota: number, cc: number, vip: number, p1: number, urgent: number }} priority_summary
 * @property {Array<ChannelData>} channels
 */

/**
 * @typedef {Object} ChannelData
 * @property {string} channel
 * @property {number} sla
 * @property {number} open
 * @property {number} closed
 * @property {Array<{ name: string, interactions: number, tickets: number, fcr: number }>} top_corporate
 * @property {Array<{ name: string, interactions: number, tickets: number, fcr: number }>} top_kip
 */
