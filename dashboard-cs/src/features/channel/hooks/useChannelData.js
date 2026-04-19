import { useState, useEffect, useCallback } from 'react';
import { fetchDashboard } from '../../../services/dashboardService';

function mapChannel(ch) {
  return {
    key: ch.channel,
    sla: ch.sla ?? 0,
    open: ch.open ?? 0,
    closed: ch.closed ?? 0,
    topCorporate: (ch.top_corporate ?? []).map((c) => ({
      name: c.company_name,
      interactions: c.interactions,
      tickets: c.tickets,
      fcr: c.fcr_percentage,
    })),
    topKip: (ch.top_kip ?? []).map((k) => ({
      name: k.topic,
      interactions: k.interactions,
      tickets: k.tickets,
      fcr: k.fcr_percentage,
    })),
  };
}

function getTodayISO() {
  return new Date().toISOString().slice(0, 10);
}

export function useChannelData(filter) {
  const today = getTodayISO();
  const from = filter?.from ?? today;
  const to = filter?.to ?? today;

  const [channels, setChannels] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await fetchDashboard({ from, to });
      setChannels((data.channels ?? []).map(mapChannel));
    } catch (err) {
      setError(err.message ?? 'Gagal memuat data channel');
    } finally {
      setLoading(false);
    }
  }, [from, to]);

  useEffect(() => {
    load();
  }, [load]);

  return { channels, loading, error, refresh: load };
}