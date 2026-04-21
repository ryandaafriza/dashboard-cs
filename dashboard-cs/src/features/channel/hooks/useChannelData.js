import { useState, useEffect, useCallback } from 'react';
import { fetchDashboard, fetchChannelDetail } from '../../../services/dashboardService';

const CHANNEL_KEYS = ['email', 'whatsapp', 'social_media', 'live_chat', 'call_center'];

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
      // Fetch summary dashboard + semua channel detail (preview page 1) secara paralel
      const [dashData, ...detailResults] = await Promise.all([
        fetchDashboard({ from, to }),
        ...CHANNEL_KEYS.map((ch) =>
          fetchChannelDetail({ from, to, channel: ch, page: 1, limit: 10 })
            .then((res) => ({ channel: ch, data: res }))
            .catch(() => ({ channel: ch, data: null })) // jangan gagal semua jika 1 channel error
        ),
      ]);

      // Buat map channelKey → detail result
      const detailMap = {};
      for (const result of detailResults) {
        detailMap[result.channel] = result.data;
      }

      const mapped = (dashData.channels ?? []).map((ch) => {
        const detail = detailMap[ch.channel];
        return {
          key: ch.channel,
          sla: ch.sla ?? 0,
          open: ch.open ?? 0,
          closed: ch.closed ?? 0,
          topCorporate: (detail?.top_corporate?.data ?? []).map((c) => ({
            name: c.company_name,
            interactions: c.interactions,
            tickets: c.tickets,
            fcr: c.fcr_percentage,
          })),
          topCorporatePagination: detail?.top_corporate?.pagination ?? null,
          topKip: (detail?.top_kip?.data ?? []).map((k) => ({
            name: k.topic,
            interactions: k.interactions,
            tickets: k.tickets,
            fcr: k.fcr_percentage,
          })),
          topKipPagination: detail?.top_kip?.pagination ?? null,
        };
      });

      setChannels(mapped);
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