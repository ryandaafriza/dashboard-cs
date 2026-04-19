import { useState, useEffect, useCallback } from 'react';
import {
  fetchActiveIncidents,
  fetchIncidentHistory,
  createIncident,
  resolveIncident,
} from '../../../services/incidentService';

// Helper: normalisasi berbagai bentuk response ke array
function toArray(data) {
  if (!data) return [];
  if (Array.isArray(data)) return data;
  if (Array.isArray(data.incidents)) return data.incidents;
  return [];
}

export function useIncidents(filter) {
  const [active,      setActive]      = useState([]);
  const [history,     setHistory]     = useState([]);
  const [loading,     setLoading]     = useState(true);
  const [histLoading, setHistLoading] = useState(false);
  const [error,       setError]       = useState(null);

  const loadActive = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await fetchActiveIncidents();
      setActive(toArray(data));
    } catch (err) {
      setError(err.message ?? 'Gagal memuat incident aktif');
      setActive([]);
    } finally {
      setLoading(false);
    }
  }, []);

  const loadHistory = useCallback(async () => {
    setHistLoading(true);
    try {
      const data = await fetchIncidentHistory({ from: filter.from, to: filter.to });
      setHistory(toArray(data));
    } catch {
      setHistory([]);
    } finally {
      setHistLoading(false);
    }
  }, [filter.from, filter.to]);

  useEffect(() => {
    loadActive();
  }, [loadActive]);

  async function addIncident(payload) {
    const result = await createIncident({
      ...payload,
      // Pastikan format datetime sesuai: "2026-04-19 09:00:00"
      started_at: new Date().toISOString().replace('T', ' ').slice(0, 19),
    });
    await loadActive(); // refresh banner setelah tambah
    return result;
  }

  async function resolve(id) {
    await resolveIncident(id);
    await loadActive(); // refresh banner setelah resolve
  }

  return {
    active, history, loading, histLoading, error,
    reload: loadActive, loadHistory, addIncident, resolve,
  };
}