import { useState, useEffect, useCallback } from "react";
import { fetchDashboard } from "../../../services/dashboardService";
import { fetchRealtime } from "../../../services/realtimeService";

function toLocaleDateRange(from, to) {
  if (!from) return "";
  const fmt = (d) =>
    new Date(d).toLocaleDateString("en-GB", {
      day: "2-digit",
      month: "short",
      year: "numeric",
    });
  return from === to ? fmt(from) : `${fmt(from)} – ${fmt(to)}`;
}

function getTodayISO() {
  return new Date().toISOString().slice(0, 10);
}

const EMPTY_DATA = {
  slaToday: 0,
  slaDelta: 0,
  slaStatus: "N/A",
  dailyTrend: [],
  ticketsPerHour: [],
  createdToday: 0,
  createdDelta: 0,
  openTickets: 0,
  unassigned: 0,
  csat: { percentage: 0, score: 0 },
  summary: { total: 0, open: 0, closed: 0 },
  priority: { roaming: 0, extra_quota: 0, cc: 0, vip: 0, p1: 0, urgent: 0 },
  lastSync: null,
  dateRange: "",
};

function mapDashboard(raw, filter) {
  const {
    summary,
    daily_trend = [],
    tickets_per_hour = [],
    priority_summary = {},
  } = raw;

  return {
    dailyTrend: daily_trend.map((d) => ({
      date: d.date, 
      created: d.created,
      solved: d.solved,
    })),
    ticketsPerHour: tickets_per_hour.map((d) => ({
      hour: d.hour.slice(0, 2),
      created: d.created,
      solved: d.solved,
    })),
    csat: {
      percentage: summary?.csat_percentage ?? 0,
      score: summary?.csat_score ?? 0,
    },
    summary: {
      total: summary?.total_tickets ?? 0,
      open: summary?.open ?? 0,
      closed: summary?.closed ?? 0,
    },
    // open, unassigned sekarang dari summary (per periode filter)
    openTickets: summary?.open ?? 0,
    unassigned: summary?.unassigned ?? 0,
    priority: {
      roaming: priority_summary.roaming ?? 0,
      extra_quota: priority_summary.extra_quota ?? 0,
      cc: priority_summary.cc ?? 0,
      vip: priority_summary.vip ?? 0,
      p1: priority_summary.p1 ?? 0,
      urgent: priority_summary.urgent ?? 0,
    },
    dateRange: toLocaleDateRange(filter?.from, filter?.to),
  };
}

function mapRealtime(rt) {
  return {
    slaToday: rt.sla_today?.percentage ?? 0,
    slaDelta: rt.sla_today?.delta ?? 0,
    slaStatus:
      (rt.sla_today?.percentage ?? 0) >= 80 ? "Achieved" : "Below Target",
    createdToday: rt.created_today?.total ?? 0,
    createdDelta: rt.created_today?.delta ?? 0,
    incidentsActive: rt.incidents_active ?? 0,
  };
}

export function useDashboardData(filter) {
  const today = getTodayISO();
  const from = filter?.from ?? today;
  const to = filter?.to ?? today;

  const [data, setData] = useState(EMPTY_DATA);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const [dashRaw, rtRaw] = await Promise.all([
        fetchDashboard({ from, to }),
        fetchRealtime(),
      ]);

      const dashMapped = mapDashboard(dashRaw, dashRaw.filter ?? { from, to });
      const rtMapped = mapRealtime(rtRaw);

      setData({
        ...EMPTY_DATA,
        ...dashMapped,
        ...rtMapped,
        lastSync: new Date().toLocaleString("id-ID"),
      });
    } catch (err) {
      setError(err.message ?? "Gagal memuat data");
    } finally {
      setLoading(false);
    }
  }, [from, to]);

  useEffect(() => {
    load();
  }, [load]);

  return { data, loading, error, refresh: load };
}
