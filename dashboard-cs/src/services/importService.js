import { apiFetch } from './apiClient';

export async function importExcel(file) {
  const formData = new FormData();
  formData.append('file', file);

  const res = await fetch(`${import.meta.env.VITE_API_BASE_URL ?? ''}/api/v1/import`, {
    method: 'POST',
    body: formData,
    // Jangan set Content-Type — browser otomatis set multipart boundary
  });

  const json = await res.json().catch(() => null);
  if (json && json.success === false) throw new Error(json.message ?? `Server error ${res.status}`);
  if (!res.ok) throw new Error(`[${res.status}] ${json?.message ?? res.statusText}`);
  return json?.data !== undefined ? json.data : json;
}

export async function syncData() {
  const res = await fetch(`${import.meta.env.VITE_API_BASE_URL ?? ''}/api/v1/sync-data`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
  });

  const json = await res.json().catch(() => null);
  if (json && json.success === false) throw new Error(json.message ?? `Server error ${res.status}`);
  if (!res.ok) throw new Error(`[${res.status}] ${json?.message ?? res.statusText}`);
  return json?.data !== undefined ? json.data : json;
}