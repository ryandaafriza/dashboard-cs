const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '';

/**
 * Generic fetch wrapper.
 * Handles response envelope: { success, message, data }
 * Throws error dengan message dari server jika success: false atau status tidak ok.
 */
export async function apiFetch(path, options = {}) {
  const url = `${BASE_URL}${path}`;

  const res = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  });

  const json = await res.json().catch(() => null);

  // Jika server mengembalikan { success: false }, lempar error dengan message dari server
  if (json && json.success === false) {
    throw new Error(json.message ?? `Server error ${res.status}`);
  }

  if (!res.ok) {
    throw new Error(`[${res.status}] ${json?.message ?? res.statusText}`);
  }

  // Unwrap envelope { success, message, data } → return data langsung
  // Jika tidak ada envelope, return json apa adanya
  return json?.data !== undefined ? json.data : json;
}

/**
 * Build query string dari plain object, skip null/undefined/empty string.
 */
export function buildQuery(params = {}) {
  const qs = Object.entries(params)
    .filter(([, v]) => v !== null && v !== undefined && v !== '')
    .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
    .join('&');
  return qs ? `?${qs}` : '';
}