import { useState, useEffect, useCallback } from 'react';
import { fetchChannelDetail } from '../../../services/dashboardService';

const DEFAULT_PAGINATION = { page: 1, limit: 10, total_items: 0, total_pages: 1 };

function mapCorporate(c) {
  return {
    name: c.company_name,
    interactions: c.interactions,
    tickets: c.tickets,
    fcr: c.fcr_percentage,
  };
}

function mapKip(k) {
  return {
    name: k.topic,
    interactions: k.interactions,
    tickets: k.tickets,
    fcr: k.fcr_percentage,
  };
}

export function useChannelDetail(filter, channelKey) {
  const today = new Date().toISOString().slice(0, 10);
  const from = filter?.from ?? today;
  const to = filter?.to ?? today;

  // Pagination state terpisah untuk corporate & kip
  const [corporatePage, setCorporatePage] = useState(1);
  const [kipPage, setKipPage] = useState(1);
  const limit = 10;

  const [corporateData, setCorporateData] = useState([]);
  const [corporatePagination, setCorporatePagination] = useState(DEFAULT_PAGINATION);
  const [kipData, setKipData] = useState([]);
  const [kipPagination, setKipPagination] = useState(DEFAULT_PAGINATION);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const load = useCallback(async () => {
    if (!channelKey) return;
    setLoading(true);
    setError(null);
    try {
      // Fetch corporate & kip page secara bersamaan
      const [corpRes, kipRes] = await Promise.all([
        fetchChannelDetail({ from, to, channel: channelKey, page: corporatePage, limit }),
        fetchChannelDetail({ from, to, channel: channelKey, page: kipPage, limit }),
      ]);

      setCorporateData((corpRes.top_corporate?.data ?? []).map(mapCorporate));
      setCorporatePagination(corpRes.top_corporate?.pagination ?? DEFAULT_PAGINATION);
      setKipData((kipRes.top_kip?.data ?? []).map(mapKip));
      setKipPagination(kipRes.top_kip?.pagination ?? DEFAULT_PAGINATION);
    } catch (err) {
      setError(err.message ?? 'Gagal memuat detail channel');
    } finally {
      setLoading(false);
    }
  }, [from, to, channelKey, corporatePage, kipPage, limit]);

  useEffect(() => {
    load();
  }, [load]);

  return {
    corporateData,
    corporatePagination,
    kipData,
    kipPagination,
    loading,
    error,
    setCorporatePage,
    setKipPage,
    refresh: load,
  };
}