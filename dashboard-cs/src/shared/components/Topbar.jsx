import React, { useState, useCallback } from 'react';
import { ImportExcelButton } from './ImportExcelButton';
import { ExportExcelButton } from './ExportExcelButton';
import { SyncButton } from './SyncButton';
import { Toast } from './Toast';

export function Topbar({ lastSync, onRefresh, loading, filter, onFilterChange }) {
  const [toast, setToast]       = useState(null);
  const [needsSync, setNeedsSync] = useState(false);

  const clearToast = useCallback(() => setToast(null), []);

  function handleFromChange(e) {
    const from = e.target.value;
    const to = filter.to < from ? from : filter.to;
    onFilterChange({ from, to });
  }

  function handleToChange(e) {
    const to = e.target.value;
    const from = filter.from > to ? to : filter.from;
    onFilterChange({ from, to });
  }

  function handleImportResult(result) {
    setToast(result);
    if (result.type === 'success') setNeedsSync(true);
  }

  async function handleSync() {
    try {
      await onRefresh();
      setNeedsSync(false);
      setToast({ type: 'success', message: 'Data berhasil disinkronisasi.' });
    } catch (err) {
      setToast({ type: 'error', message: err.message ?? 'Gagal sinkronisasi data.' });
    }
  }

  return (
    <>
      <header className="topbar">
        <div className="topbar__left">
          <h1 className="topbar__title">Dashboard</h1>
          <span className="topbar__badge">Overview</span>
        </div>

        <div className="topbar__right">
          <div className="sync-info">
            <div className="sync-dot" />
            <span>Last Sync: {lastSync ?? '—'}</span>
          </div>

          <div className="date-filter">
            <label className="date-filter__label">From</label>
            <input
              type="date"
              className="date-filter__input"
              value={filter.from}
              max={filter.to}
              onChange={handleFromChange}
              disabled={loading}
            />
            <span className="date-filter__sep">–</span>
            <label className="date-filter__label">To</label>
            <input
              type="date"
              className="date-filter__input"
              value={filter.to}
              min={filter.from}
              max={new Date().toISOString().slice(0, 10)}
              onChange={handleToChange}
              disabled={loading}
            />
          </div>

          {/* Sync — selalu ada, warna berubah sesuai state */}
          <SyncButton onSync={handleSync} needsSync={needsSync} />

          <ImportExcelButton onImportSuccess={handleImportResult} />

          <ExportExcelButton filter={filter} onResult={setToast} />
        </div>
      </header>

      <Toast toast={toast} onClose={clearToast} />

      <style>{`
        .date-filter {
          display: flex; align-items: center; gap: 6px;
        }
        .date-filter__label { font-size: 12px; color: var(--text-muted, #888); }
        .date-filter__input {
          font-size: 13px; padding: 4px 8px;
          border: 1px solid var(--border-default, #ddd);
          border-radius: var(--radius-sm, 6px);
          background: var(--bg-card, #fff);
          color: var(--text-primary, #333);
          cursor: pointer;
        }
        .date-filter__input:disabled { opacity: 0.5; cursor: not-allowed; }
        .date-filter__sep { color: var(--text-muted, #888); }

        .topbar-btn {
          display: flex; align-items: center; gap: 6px;
          padding: 6px 12px; font-size: 13px; font-weight: 500;
          border-radius: var(--radius-sm, 6px);
          border: 1px solid transparent;
          cursor: pointer; transition: opacity 0.2s;
        }
        .topbar-btn:disabled { opacity: 0.6; cursor: not-allowed; }
        .topbar-btn--import {
          background: var(--color-info-bg, rgba(59,130,246,0.1));
          border-color: var(--color-info-border, rgba(59,130,246,0.25));
          color: var(--color-info, #3b82f6);
        }
        .topbar-btn--import:hover:not(:disabled) { opacity: 0.8; }
      `}</style>
    </>
  );
}