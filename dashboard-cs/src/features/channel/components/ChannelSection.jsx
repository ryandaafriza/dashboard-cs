import React from 'react';
import { ChannelCard } from './ChannelCard';
import { useChannelData } from '../hooks/useChannelData';

export function ChannelSection({ filter }) {
  const { channels, loading, error } = useChannelData(filter);

  if (loading) return <div className="section-loading">Memuat data channel…</div>;
  if (error) return <div className="section-error">Error: {error}</div>;

  return (
    <section>
      <div className="section-label">Channel Performance</div>
      <div className="channel-section-grid">
        {channels.map((ch) => (
          <ChannelCard key={ch.key} channel={ch} filter={filter} />
        ))}
      </div>
    </section>
  );
}