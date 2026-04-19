import React from 'react';
import { SLACard } from './SLACard';
import { DailyTrendCard } from './DailyTrendCard';
import { TicketsPerHourCard } from './TicketsPerHourCard';
import { CreatedTodayCard, OpenTicketsCard, UnassignedCard, CSATCard } from './StatCards';

export function KPISection({ data }) {
  return (
    <section>
      <div className="section-label">Key Performance Indicators</div>
      <div className="kpi-section">
        <SLACard
          slaToday={data.slaToday}
          slaDelta={data.slaDelta}
          slaStatus={data.slaStatus}
        />
        <DailyTrendCard data={data.dailyTrend} />
        <TicketsPerHourCard data={data.ticketsPerHour} />
        <CreatedTodayCard value={data.createdToday} delta={data.createdDelta} />
        <OpenTicketsCard value={data.openTickets} />
        <UnassignedCard value={data.unassigned} />
        <CSATCard
          percentage={data.csat.percentage}
          score={data.csat.score}
        />
      </div>
    </section>
  );
}
