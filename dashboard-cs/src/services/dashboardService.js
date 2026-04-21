import { apiFetch, buildQuery } from './apiClient';

/**
 * Fetch filtered dashboard data.
 * @param {{ from: string, to: string }} filter
 */
export async function fetchDashboard({ from, to }) {
  const qs = buildQuery({ from, to });
  return apiFetch(`/api/v1/dashboard${qs}`);
}

/**
 * Fetch top_corporate & top_kip per channel dengan pagination.
 * @param {{ from: string, to: string, channel: string, page?: number, limit?: number }} params
 */
export async function fetchChannelDetail({ from, to, channel, page = 1, limit = 10 }) {
  const qs = buildQuery({ from, to, channel, page, limit });
  return apiFetch(`/api/v1/dashboard/channels${qs}`);
}

/**
 * @typedef {Object} DashboardResponse
 * @property {{ from: string, to: string }} filter
 * @property {{ total_tickets: number, open: number, closed: number, csat_percentage: number, csat_score: number }} summary
 * @property {Array<{ date: string, created: number, solved: number }>} daily_trend
 * @property {Array<{ hour: string, created: number, solved: number }>} tickets_per_hour
 * @property {{ roaming: number, extra_quota: number, cc: number, vip: number, p1: number, urgent: number }} priority_summary
 * @property {Array<ChannelSummary>} channels
 */

/**
 * @typedef {Object} ChannelSummary
 * @property {string} channel
 * @property {number} sla
 * @property {number} open
 * @property {number} closed
 */

/**
 * @typedef {Object} ChannelDetailResponse
 * @property {{ data: Array, pagination: { page: number, limit: number, total_items: number, total_pages: number } }} top_corporate
 * @property {{ data: Array, pagination: { page: number, limit: number, total_items: number, total_pages: number } }} top_kip
 */