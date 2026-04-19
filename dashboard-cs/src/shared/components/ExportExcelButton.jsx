import React, { useState } from 'react';
import { exportExcel } from '../../services/exportService';

const CHANNEL_OPTIONS = [
  { value: 'all',          label: 'All Channel' },
  { value: 'email',        label: 'Email' },
  { value: 'whatsapp',     label: 'WhatsApp' },
  { value: 'social_media', label: 'Social Media' },
  { value: 'live_chat',    label: 'Live Chat' },
  { value: 'call_center',  label: 'Call Center 188' },
];

export function ExportExcelButton({ filter, onResult }) {
  const [open, setOpen]       = useState(false);
  const [selected, setSelected] = useState(['all']);
  const [loading, setLoading] = useState(false);

  function toggleChannel(value) {
    if (value === 'all') {
      setSelected(['all']);
      return;
    }
    setSelected((prev) => {
      const withoutAll = prev.filter((v) => v !== 'all');
      if (withoutAll.includes(value)) {
        const next = withoutAll.filter((v) => v !== value);
        return next.length === 0 ? ['all'] : next;
      }
      return [...withoutAll, value];
    });
  }

  async function handleExport() {
    setLoading(true);
    try {
      const channel = selected.includes('all') ? 'all' : selected.join(',');
      await exportExcel({ from: filter.from, to: filter.to, channel });
      setOpen(false);
      onResult({ type: 'success', message: 'File Excel berhasil didownload.' });
    } catch (err) {
      onResult({ type: 'error', message: err.message ?? 'Gagal export data.' });
    } finally {
      setLoading(false);
    }
  }

  const channelLabel = selected.includes('all')
    ? 'All Channel'
    : selected.map((v) => CHANNEL_OPTIONS.find((o) => o.value === v)?.label).join(', ');

  return (
    <>
      <button
        className="topbar-btn topbar-btn--export"
        onClick={() => setOpen(true)}
        title="Export ke Excel"
      >
        <span>⬇</span>
        Export Excel
      </button>

      {/* Modal overlay */}
      {open && (
        <div className="export-overlay" onClick={() => !loading && setOpen(false)}>
          <div className="export-modal" onClick={(e) => e.stopPropagation()}>
            <div className="export-modal__header">
              <span className="export-modal__title">Export Excel</span>
              <button className="export-modal__close" onClick={() => !loading && setOpen(false)}>×</button>
            </div>

            <div className="export-modal__body">
              {/* Info periode */}
              <div className="export-modal__info">
                <span className="export-modal__info-label">Periode</span>
                <span className="export-modal__info-value">{filter.from} s/d {filter.to}</span>
              </div>

              {/* Pilih channel */}
              <div className="export-modal__section-label">Filter Channel</div>
              <div className="export-channel-grid">
                {CHANNEL_OPTIONS.map((opt) => {
                  const isActive = selected.includes(opt.value) ||
                    (opt.value !== 'all' && selected.includes('all'));
                  const isAll    = opt.value === 'all';
                  return (
                    <button
                      key={opt.value}
                      className={`export-channel-btn ${isActive ? 'export-channel-btn--active' : ''} ${isAll ? 'export-channel-btn--all' : ''}`}
                      onClick={() => toggleChannel(opt.value)}
                      disabled={loading}
                    >
                      {isActive && <span className="export-channel-btn__check">✓</span>}
                      {opt.label}
                    </button>
                  );
                })}
              </div>

              <div className="export-modal__preview">
                Akan export: <strong>{channelLabel}</strong>
              </div>
            </div>

            <div className="export-modal__footer">
              <button
                className="export-modal__cancel"
                onClick={() => setOpen(false)}
                disabled={loading}
              >
                Batal
              </button>
              <button
                className="export-modal__submit"
                onClick={handleExport}
                disabled={loading}
              >
                {loading
                  ? <><span className="spin">↻</span> Mengexport…</>
                  : <><span>⬇</span> Download Excel</>
                }
              </button>
            </div>
          </div>
        </div>
      )}

      <style>{`
        .topbar-btn--export {
          background: var(--color-success-bg, rgba(34,197,94,0.1));
          border-color: var(--color-success-border, rgba(34,197,94,0.25));
          color: var(--color-success, #22c55e);
        }
        .topbar-btn--export:hover { opacity: 0.8; }

        .export-overlay {
          position: fixed;
          inset: 0;
          z-index: 1000;
          background: rgba(0,0,0,0.55);
          display: flex;
          align-items: center;
          justify-content: center;
          animation: fadeIn 0.15s ease;
        }
        .export-modal {
          background: var(--bg-card, #1c2333);
          border: 1px solid var(--border-default, rgba(255,255,255,0.08));
          border-radius: var(--radius-lg, 14px);
          width: 420px;
          max-width: calc(100vw - 32px);
          box-shadow: var(--shadow-elevated, 0 4px 24px rgba(0,0,0,0.5));
          animation: slideUp 0.18s ease;
        }
        .export-modal__header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 16px 20px 12px;
          border-bottom: 1px solid var(--border-subtle, rgba(255,255,255,0.05));
        }
        .export-modal__title {
          font-size: 15px;
          font-weight: 600;
          color: var(--text-primary, #e8edf5);
        }
        .export-modal__close {
          background: none;
          border: none;
          color: var(--text-muted, #4d5a70);
          font-size: 20px;
          cursor: pointer;
          line-height: 1;
          padding: 0 2px;
        }
        .export-modal__close:hover { color: var(--text-primary, #e8edf5); }

        .export-modal__body {
          padding: 16px 20px;
          display: flex;
          flex-direction: column;
          gap: 14px;
        }
        .export-modal__info {
          display: flex;
          align-items: center;
          justify-content: space-between;
          background: var(--bg-surface, #161b26);
          border-radius: var(--radius-sm, 6px);
          padding: 8px 12px;
          font-size: 13px;
        }
        .export-modal__info-label { color: var(--text-muted, #4d5a70); }
        .export-modal__info-value { color: var(--text-primary, #e8edf5); font-weight: 500; }

        .export-modal__section-label {
          font-size: 12px;
          font-weight: 500;
          color: var(--text-muted, #4d5a70);
          text-transform: uppercase;
          letter-spacing: 0.5px;
        }
        .export-channel-grid {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 8px;
        }
        .export-channel-btn {
          display: flex;
          align-items: center;
          gap: 6px;
          padding: 8px 12px;
          background: var(--bg-surface, #161b26);
          border: 1px solid var(--border-default, rgba(255,255,255,0.08));
          border-radius: var(--radius-sm, 6px);
          color: var(--text-secondary, #8a95a8);
          font-size: 13px;
          cursor: pointer;
          transition: all 0.15s;
          text-align: left;
        }
        .export-channel-btn--all { grid-column: 1 / -1; }
        .export-channel-btn:hover:not(:disabled) {
          border-color: var(--border-strong, rgba(255,255,255,0.14));
          color: var(--text-primary, #e8edf5);
        }
        .export-channel-btn--active {
          background: var(--color-info-bg, rgba(59,130,246,0.1));
          border-color: var(--color-info-border, rgba(59,130,246,0.3));
          color: var(--color-info, #3b82f6);
        }
        .export-channel-btn__check { font-size: 11px; }
        .export-channel-btn:disabled { opacity: 0.5; cursor: not-allowed; }

        .export-modal__preview {
          font-size: 12px;
          color: var(--text-muted, #4d5a70);
        }
        .export-modal__preview strong {
          color: var(--text-secondary, #8a95a8);
        }

        .export-modal__footer {
          display: flex;
          align-items: center;
          justify-content: flex-end;
          gap: 8px;
          padding: 12px 20px 16px;
          border-top: 1px solid var(--border-subtle, rgba(255,255,255,0.05));
        }
        .export-modal__cancel {
          padding: 7px 16px;
          font-size: 13px;
          background: none;
          border: 1px solid var(--border-default, rgba(255,255,255,0.08));
          border-radius: var(--radius-sm, 6px);
          color: var(--text-secondary, #8a95a8);
          cursor: pointer;
        }
        .export-modal__cancel:hover:not(:disabled) {
          border-color: var(--border-strong);
          color: var(--text-primary);
        }
        .export-modal__submit {
          display: flex;
          align-items: center;
          gap: 6px;
          padding: 7px 16px;
          font-size: 13px;
          font-weight: 500;
          background: var(--color-success-bg, rgba(34,197,94,0.12));
          border: 1px solid var(--color-success-border, rgba(34,197,94,0.3));
          border-radius: var(--radius-sm, 6px);
          color: var(--color-success, #22c55e);
          cursor: pointer;
        }
        .export-modal__submit:hover:not(:disabled) { opacity: 0.8; }
        .export-modal__submit:disabled { opacity: 0.6; cursor: not-allowed; }

        .spin { display: inline-block; animation: spin 0.8s linear infinite; }
        @keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
        @keyframes slideUp {
          from { opacity: 0; transform: translateY(10px); }
          to   { opacity: 1; transform: translateY(0); }
        }
      `}</style>
    </>
  );
}