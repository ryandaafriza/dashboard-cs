import { buildQuery } from './apiClient';

/**
 * Download export Excel dari server sebagai file.
 * @param {{ from: string, to: string, channel: string }} params
 */
export async function exportExcel({ from, to, channel = 'all' }) {
  const qs = buildQuery({ from, to, channel });
  const url = `${import.meta.env.VITE_API_BASE_URL ?? ''}/api/v1/export${qs}`;

  const res = await fetch(url, { method: 'GET' });

  if (!res.ok) {
    // Coba baca pesan error dari body
    const json = await res.json().catch(() => null);
    throw new Error(json?.message ?? `Export gagal [${res.status}]`);
  }

  // Ambil nama file dari Content-Disposition header jika ada
  const disposition = res.headers.get('Content-Disposition') ?? '';
  const match = disposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/);
  const filename = match ? match[1].replace(/['"]/g, '') : `Dashboard_Report_${from}_${to}_${channel}.xlsx`;

  const blob = await res.blob();
  const objectUrl = URL.createObjectURL(blob);

  // Trigger download
  const a = document.createElement('a');
  a.href = objectUrl;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  a.remove();
  URL.revokeObjectURL(objectUrl);
}