import { apiFetch, buildQuery } from './apiClient';

export async function fetchActiveIncidents() {
  // Response setelah unwrap envelope: { count: number, incidents: [...] }
  return apiFetch('/api/v1/incidents/active');
}

export async function fetchIncidentHistory({ from, to }) {
  const qs = buildQuery({ from, to });
  // Response setelah unwrap envelope: { incidents: [...] } atau array langsung
  return apiFetch(`/api/v1/incidents/history${qs}`);
}

export async function createIncident({ title, description, severity, started_at, created_by }) {
  return apiFetch('/api/v1/incidents', {
    method: 'POST',
    body: JSON.stringify({ title, description, severity, started_at, created_by }),
  });
}

export async function resolveIncident(id) {
  return apiFetch(`/api/v1/incidents/${id}/resolve`, { method: 'PATCH' });
}