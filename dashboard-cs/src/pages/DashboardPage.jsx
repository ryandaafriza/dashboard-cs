import React, { useState } from 'react';
import { Topbar } from '../shared/components/Topbar';
import { KPISection } from '../features/kpi/components/KPISection';
import { SummarySection } from '../features/summary/components/SummarySection';
import { PrioritySection } from '../features/priority/components/PrioritySection';
import { ChannelSection } from '../features/channel/components/ChannelSection';
import { IncidentSection } from '../features/incident/components/IncidentSection';
import { useDashboardData } from '../features/kpi/hooks/useDashboardData';

function getTodayISO() {
  return new Date().toISOString().slice(0, 10);
}

export function DashboardPage() {
  const today = getTodayISO();
  const [filter, setFilter] = useState({ from: today, to: today });

  const { data, loading, error, refresh } = useDashboardData(filter);

  return (
    <div className="dashboard-page">
      <Topbar
        lastSync={data.lastSync}
        onRefresh={refresh}
        loading={loading}
        filter={filter}
        onFilterChange={setFilter}
      />
      {error && (
        <div className="dashboard-error-banner">⚠️ {error}</div>
      )}
      <main className="dashboard-content">
        <KPISection data={data} />
        <SummarySection summary={data.summary} />
        <PrioritySection priority={data.priority} />
        <ChannelSection filter={filter} />
        <IncidentSection filter={filter} />
      </main>
    </div>
  );
}