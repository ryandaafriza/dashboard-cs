import React, { useState } from 'react';
import { CHANNEL_CONFIG } from '../utils/channelUtils';
import { ChannelMetrics } from './ChannelMetrics';
import { useChannelDetail } from '../hooks/useChannelDetail';

/* ─── Pagination ──────────────────────────────────────────────────── */
function Pagination({ pagination, onPageChange, disabled }) {
  const { page, total_pages } = pagination;
  if (total_pages <= 1) return null;
  return (
    <div className="ch-pagination">
      <button
        className="ch-pagination__btn"
        onClick={() => onPageChange(page - 1)}
        disabled={disabled || page <= 1}
      >‹</button>
      <span className="ch-pagination__info">{page} / {total_pages}</span>
      <button
        className="ch-pagination__btn"
        onClick={() => onPageChange(page + 1)}
        disabled={disabled || page >= total_pages}
      >›</button>
    </div>
  );
}

/* ─── Preview Table (tampil langsung di card dashboard) ───────────── */
function PreviewTable({ title, rows, color, totalItems, onShowAll }) {
  return (
    <div className="channel-table">
      <div className="channel-table__title-row">
        <div className="channel-table__title" style={{ borderLeft: `2px solid ${color}` }}>
          {title}
        </div>
        {totalItems > 0 && (
          <button className="channel-table__more-btn" onClick={onShowAll}>
            → Lihat Semua ({totalItems})
          </button>
        )}
      </div>

      {!rows.length ? (
        <div className="channel-table__empty">No data available</div>
      ) : (
        <table className="channel-table__grid">
          <thead>
            <tr>
              <th className="channel-table__th">Nama</th>
              <th className="channel-table__th channel-table__th--num">Interaksi</th>
              <th className="channel-table__th channel-table__th--num">Tiket</th>
              <th className="channel-table__th channel-table__th--num">%FCR</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row, i) => (
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
    </div>
  );
}

/* ─── Modal dengan Tabs + Pagination ─────────────────────────────── */
function ChannelDetailModal({ channelKey, config, filter, onClose }) {
  // const [activeTab, setActiveTab] = useState('corporate'); // reserved for KIP tab later
  const {
    corporateData, corporatePagination,
    // kipData, kipPagination,   // reserved for KIP tab later
    loading, error,
    setCorporatePage,
    // setKipPage,               // reserved for KIP tab later
  } = useChannelDetail(filter, channelKey);

  return (
    <div className="ch-overlay" onClick={onClose}>
      <div className="ch-modal" onClick={(e) => e.stopPropagation()}>

        <div className="ch-modal__header">
          <span className="ch-modal__title">
            <span style={{ marginRight: 6 }}>{config.icon}</span>
            {config.label} — Top Consument
          </span>
          <button className="ch-modal__close" onClick={onClose}>×</button>
        </div>

        {/* Tabs — KIP hidden for now */}
        {/* <div className="ch-modal__tabs">
          <button className="ch-modal__tab ch-modal__tab--active" style={{ '--tab-color': config.color }}>
            Top Corporate
            {corporatePagination.total_items > 0 && <span className="ch-modal__tab-badge">{corporatePagination.total_items}</span>}
          </button>
          <button
            className={`ch-modal__tab ${activeTab === 'kip' ? 'ch-modal__tab--active' : ''}`}
            onClick={() => setActiveTab('kip')}
            style={{ '--tab-color': config.color }}
          >
            Top KIP
            {kipPagination.total_items > 0 && <span className="ch-modal__tab-badge">{kipPagination.total_items}</span>}
          </button>
        </div> */}

        <div className="ch-modal__body">
          {loading ? (
            <div className="ch-modal__state">Memuat data…</div>
          ) : error ? (
            <div className="ch-modal__state ch-modal__state--error">{error}</div>
          ) : !corporateData.length ? (
            <div className="ch-modal__state">Tidak ada data.</div>
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
                {corporateData.map((row, i) => {
                  const rowNum = (corporatePagination.page - 1) * (corporatePagination.limit ?? 10) + i + 1;
                  return (
                    <tr key={i}>
                      <td className="ch-modal__num">{rowNum}</td>
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
                  );
                })}
              </tbody>
            </table>
          )}
        </div>

        {!loading && !error && corporateData.length > 0 && (
          <div className="ch-modal__footer">
            <span className="ch-modal__total">Total: {corporatePagination.total_items} data</span>
            <Pagination pagination={corporatePagination} onPageChange={setCorporatePage} disabled={loading} />
          </div>
        )}
      </div>

      <style>{modalStyles}</style>
    </div>
  );
}

/* ─── ChannelCard ─────────────────────────────────────────────────── */
export function ChannelCard({ channel, filter }) {
  // null = tutup, 'corporate' / 'kip' = buka modal di tab tsb
  const [modalTab, setModalTab] = useState(null);

  const config = CHANNEL_CONFIG[channel.key] ?? {
    label: channel.key, icon: '•', color: '#8a95a8', gradient: 'transparent',
  };

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

          <div className="channel-card__tables">
            {/* Top Corporate preview */}
            <PreviewTable
              title="Top Consument"
              rows={channel.topCorporate ?? []}
              color={config.color}
              totalItems={channel.topCorporatePagination?.total_items ?? 0}
              onShowAll={() => setModalTab('corporate')}
            />

            <div className="channel-card__divider" />

            {/* Top KIP preview */}
            {/* <PreviewTable
              title="Top KIP"
              rows={channel.topKip ?? []}
              color={config.color}
              totalItems={channel.topKipPagination?.total_items ?? 0}
              onShowAll={() => setModalTab('kip')}
            /> */}
          </div>

        </div>
      </div>

      {modalTab && (
        <ChannelDetailModal
          channelKey={channel.key}
          config={config}
          filter={filter}
          defaultTab={modalTab}
          onClose={() => setModalTab(null)}
        />
      )}

      <style>{cardStyles}</style>
    </>
  );
}

/* ─── Styles ──────────────────────────────────────────────────────── */
const cardStyles = `
  .channel-table__title-row {
    display: flex; align-items: center;
    justify-content: space-between; margin-bottom: 8px;
  }
  .channel-table__title-row .channel-table__title { margin-bottom: 0; }
  .channel-table__more-btn {
    font-size: 10px; color: var(--accent-primary, #3b82f6);
    background: none; border: none; cursor: pointer;
    padding: 2px 4px; border-radius: 4px; transition: opacity 0.15s;
  }
  .channel-table__more-btn:hover { opacity: 0.7; }
`;

const modalStyles = `
  .ch-overlay {
    position: fixed; inset: 0; z-index: 1000;
    background: rgba(0,0,0,0.6);
    display: flex; align-items: center; justify-content: center;
    animation: fadeIn 0.15s ease;
  }
  .ch-modal {
    background: var(--bg-card); border: 1px solid var(--border-default);
    border-radius: var(--radius-lg); width: 600px;
    max-width: calc(100vw - 32px); max-height: 85vh;
    display: flex; flex-direction: column;
    box-shadow: var(--shadow-elevated); animation: slideUp 0.18s ease; overflow: hidden;
  }
  .ch-modal__header {
    display: flex; align-items: center; justify-content: space-between;
    padding: 14px 20px 12px; border-bottom: 1px solid var(--border-subtle); flex-shrink: 0;
  }
  .ch-modal__title { font-size: 14px; font-weight: 600; color: var(--text-primary); display: flex; align-items: center; }
  .ch-modal__close { background: none; border: none; color: var(--text-muted); font-size: 20px; cursor: pointer; line-height: 1; }
  .ch-modal__close:hover { color: var(--text-primary); }
  .ch-modal__tabs { display: flex; border-bottom: 1px solid var(--border-subtle); flex-shrink: 0; padding: 0 20px; }
  .ch-modal__tab {
    position: relative; background: none; border: none;
    padding: 10px 14px; font-size: 13px; font-weight: 500;
    color: var(--text-muted); cursor: pointer;
    display: flex; align-items: center; gap: 6px;
    border-bottom: 2px solid transparent; margin-bottom: -1px;
    transition: color 0.15s, border-color 0.15s;
  }
  .ch-modal__tab:hover { color: var(--text-primary); }
  .ch-modal__tab--active { color: var(--tab-color); border-bottom-color: var(--tab-color); }
  .ch-modal__tab-badge {
    background: rgba(255,255,255,0.08); border-radius: 10px;
    padding: 1px 6px; font-size: 10px; font-family: var(--font-mono);
  }
  .ch-modal__body { padding: 16px 20px; overflow-y: auto; flex: 1; }
  .ch-modal__state { text-align: center; padding: 32px; color: var(--text-muted); font-size: 13px; }
  .ch-modal__state--error { color: #ef4444; }
  .ch-modal__table { width: 100%; border-collapse: collapse; font-size: 13px; }
  .ch-modal__table th {
    text-align: left; font-size: 11px; font-weight: 500;
    color: var(--text-muted); text-transform: uppercase;
    letter-spacing: 0.4px; padding: 6px 10px; border-bottom: 1px solid var(--border-subtle);
  }
  .ch-modal__table td { padding: 9px 10px; border-bottom: 1px solid rgba(255,255,255,0.03); color: var(--text-secondary); }
  .ch-modal__table tr:last-child td { border-bottom: none; }
  .ch-modal__table tr:hover td { background: rgba(255,255,255,0.02); color: var(--text-primary); }
  .ch-modal__num { text-align: right; font-family: var(--font-mono); font-size: 12px; }
  .ch-modal__footer {
    display: flex; align-items: center; justify-content: space-between;
    padding: 10px 20px; border-top: 1px solid var(--border-subtle); flex-shrink: 0;
  }
  .ch-modal__total { font-size: 11px; color: var(--text-muted); font-family: var(--font-mono); }
  .ch-pagination { display: flex; align-items: center; gap: 8px; }
  .ch-pagination__btn {
    background: rgba(255,255,255,0.05); border: 1px solid var(--border-subtle);
    border-radius: var(--radius-sm); color: var(--text-secondary);
    width: 28px; height: 28px; font-size: 14px; cursor: pointer;
    display: flex; align-items: center; justify-content: center;
    transition: background 0.15s, color 0.15s;
  }
  .ch-pagination__btn:hover:not(:disabled) { background: rgba(255,255,255,0.1); color: var(--text-primary); }
  .ch-pagination__btn:disabled { opacity: 0.3; cursor: not-allowed; }
  .ch-pagination__info { font-size: 12px; color: var(--text-muted); font-family: var(--font-mono); min-width: 40px; text-align: center; }
  @keyframes fadeIn { from{opacity:0} to{opacity:1} }
  @keyframes slideUp { from{opacity:0;transform:translateY(10px)} to{opacity:1;transform:translateY(0)} }
`;