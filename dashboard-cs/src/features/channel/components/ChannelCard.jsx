import React, { useState } from 'react';
import { CHANNEL_CONFIG } from '../utils/channelUtils';
import { ChannelMetrics } from './ChannelMetrics';
import { ChannelTable } from './ChannelTable';

function CorporateModal({ channel, config, onClose }) {
  return (
    <div className="ch-overlay" onClick={onClose}>
      <div className="ch-modal" onClick={(e) => e.stopPropagation()}>
        <div className="ch-modal__header">
          <span className="ch-modal__title">
            <span style={{ marginRight: 6 }}>{config.icon}</span>
            {config.label} — Top Corporate
          </span>
          <button className="ch-modal__close" onClick={onClose}>×</button>
        </div>
        <div className="ch-modal__body">
          {!channel.topCorporate?.length ? (
            <div className="ch-modal__empty">Tidak ada data.</div>
          ) : (
            <table className="ch-modal__table">
              <thead>
                <tr>
                  <th>#</th>
                  <th>Nama</th>
                  <th>Interaksi</th>
                  <th>Tiket</th>
                  <th>%FCR</th>
                </tr>
              </thead>
              <tbody>
                {channel.topCorporate.map((row, i) => (
                  <tr key={i}>
                    <td className="ch-modal__num">{i + 1}</td>
                    <td>{row.name}</td>
                    <td className="ch-modal__num">{row.interactions}</td>
                    <td className="ch-modal__num">{row.tickets}</td>
                    <td className="ch-modal__num" style={{
                      color: row.fcr >= 80 ? '#22c55e' : row.fcr >= 50 ? '#f59e0b' : 'var(--text-muted)',
                      fontWeight: 600,
                    }}>
                      {row.fcr > 0 ? `${row.fcr.toFixed(1)}%` : '—'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>

      <style>{`
        .ch-overlay {
          position: fixed; inset: 0; z-index: 1000;
          background: rgba(0,0,0,0.6);
          display: flex; align-items: center; justify-content: center;
          animation: fadeIn 0.15s ease;
        }
        .ch-modal {
          background: var(--bg-card);
          border: 1px solid var(--border-default);
          border-radius: var(--radius-lg);
          width: 560px;
          max-width: calc(100vw - 32px);
          max-height: 80vh;
          display: flex; flex-direction: column;
          box-shadow: var(--shadow-elevated);
          animation: slideUp 0.18s ease;
          overflow: hidden;
        }
        .ch-modal__header {
          display: flex; align-items: center; justify-content: space-between;
          padding: 14px 20px 12px;
          border-bottom: 1px solid var(--border-subtle);
          flex-shrink: 0;
        }
        .ch-modal__title {
          font-size: 14px; font-weight: 600;
          color: var(--text-primary);
          display: flex; align-items: center;
        }
        .ch-modal__close {
          background: none; border: none;
          color: var(--text-muted); font-size: 20px;
          cursor: pointer; line-height: 1;
        }
        .ch-modal__close:hover { color: var(--text-primary); }
        .ch-modal__body {
          padding: 16px 20px;
          overflow-y: auto; flex: 1;
        }
        .ch-modal__empty {
          text-align: center; padding: 32px;
          color: var(--text-muted); font-size: 13px;
        }
        .ch-modal__table {
          width: 100%; border-collapse: collapse; font-size: 13px;
        }
        .ch-modal__table th {
          text-align: left; font-size: 11px; font-weight: 500;
          color: var(--text-muted); text-transform: uppercase;
          letter-spacing: 0.4px; padding: 6px 10px;
          border-bottom: 1px solid var(--border-subtle);
        }
        .ch-modal__table td {
          padding: 9px 10px;
          border-bottom: 1px solid rgba(255,255,255,0.03);
          color: var(--text-secondary);
        }
        .ch-modal__table tr:last-child td { border-bottom: none; }
        .ch-modal__table tr:hover td { background: rgba(255,255,255,0.02); color: var(--text-primary); }
        .ch-modal__num { text-align: right; font-family: var(--font-mono); font-size: 12px; }
        @keyframes fadeIn { from{opacity:0} to{opacity:1} }
        @keyframes slideUp {
          from{opacity:0;transform:translateY(10px)}
          to{opacity:1;transform:translateY(0)}
        }
      `}</style>
    </div>
  );
}

export function ChannelCard({ channel }) {
  const [showModal, setShowModal] = useState(false);

  const config = CHANNEL_CONFIG[channel.key] ?? {
    label: channel.key, icon: '•', color: '#8a95a8', gradient: 'transparent',
  };

  const preview = channel.topCorporate?.slice(0, 10) ?? [];
  const hasMore = (channel.topCorporate?.length ?? 0) > 10;

  return (
    <>
      <div className="card channel-card animate-in" style={{ '--ch-color': config.color }}>
        <div className="card__inner channel-card__inner" style={{ background: config.gradient }}>

          <div className="channel-card__header">
            <div className="channel-card__title">
              <span className="channel-card__icon">{config.icon}</span>
              <span className="channel-card__label">{config.label}</span>
            </div>
          </div>

          <ChannelMetrics sla={channel.sla} open={channel.open} closed={channel.closed} />

          <div className="channel-card__divider" />

          {/* Hanya Top Corporate, Top KIP di-hide */}
          <div className="channel-card__tables">
            <div className="channel-table">
              <div className="channel-table__title-row">
                <div
                  className="channel-table__title"
                  style={{ borderLeft: `2px solid ${config.color}` }}
                >
                  Top Corporate
                </div>
                {(channel.topCorporate?.length > 0) && (
                  <button
                    className="channel-table__more-btn"
                    onClick={() => setShowModal(true)}
                    title="Lihat semua data"
                  >
                    → Lihat Semua ({channel.topCorporate.length})
                  </button>
                )}
              </div>

              {!preview.length ? (
                <div className="channel-table__empty">No data available</div>
              ) : (
                <table className="channel-table__grid">
                  <thead>
                    <tr>
                      <th className="channel-table__th">Name</th>
                      <th className="channel-table__th channel-table__th--num">Interaksi</th>
                      <th className="channel-table__th channel-table__th--num">Tiket</th>
                      <th className="channel-table__th channel-table__th--num">%FCR</th>
                    </tr>
                  </thead>
                  <tbody>
                    {preview.map((row, i) => (
                      <tr key={i} className="channel-table__row">
                        <td className="channel-table__td channel-table__td--name" title={row.name}>
                          {row.name}
                        </td>
                        <td className="channel-table__td channel-table__td--num">{row.interactions}</td>
                        <td className="channel-table__td channel-table__td--num">{row.tickets}</td>
                        <td className="channel-table__td channel-table__td--num">
                          <span className="channel-table__fcr" style={{
                            color: row.fcr >= 80 ? '#22c55e' : row.fcr >= 50 ? '#f59e0b' : 'var(--text-muted)',
                          }}>
                            {row.fcr > 0 ? row.fcr.toFixed(2) : '—'}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}

              {hasMore && (
                <button
                  className="channel-table__show-more"
                  onClick={() => setShowModal(true)}
                >
                  +{channel.topCorporate.length - 10} data lainnya →
                </button>
              )}
            </div>
          </div>

        </div>
      </div>

      {showModal && (
        <CorporateModal
          channel={channel}
          config={config}
          onClose={() => setShowModal(false)}
        />
      )}

      <style>{`
        .channel-table__title-row {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 8px;
        }
        .channel-table__title-row .channel-table__title {
          margin-bottom: 0;
        }
        .channel-table__more-btn {
          font-size: 10px;
          color: var(--accent-primary, #3b82f6);
          background: none;
          border: none;
          cursor: pointer;
          padding: 2px 4px;
          border-radius: 4px;
          transition: opacity 0.15s;
        }
        .channel-table__more-btn:hover { opacity: 0.7; }
        .channel-table__show-more {
          width: 100%;
          margin-top: 6px;
          padding: 5px;
          font-size: 11px;
          color: var(--accent-primary, #3b82f6);
          background: rgba(59,130,246,0.05);
          border: 1px dashed rgba(59,130,246,0.2);
          border-radius: var(--radius-sm);
          cursor: pointer;
          text-align: center;
          transition: background 0.15s;
        }
        .channel-table__show-more:hover { background: rgba(59,130,246,0.1); }
      `}</style>
    </>
  );
}