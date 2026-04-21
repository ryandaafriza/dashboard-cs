import React, { useState } from 'react';
import { useIncidents } from '../hooks/useIncidents';

const SEVERITY_CONFIG = {
  critical: { label: 'Critical', color: '#ef4444', bg: 'rgba(239,68,68,0.12)', border: 'rgba(239,68,68,0.3)' },
  high:     { label: 'High',     color: '#f97316', bg: 'rgba(249,115,22,0.12)', border: 'rgba(249,115,22,0.3)' },
  medium:   { label: 'Medium',   color: '#f59e0b', bg: 'rgba(245,158,11,0.12)', border: 'rgba(245,158,11,0.3)' },
  low:      { label: 'Low',      color: '#3b82f6', bg: 'rgba(59,130,246,0.12)', border: 'rgba(59,130,246,0.3)' },
};

function SeverityBadge({ severity }) {
  const cfg = SEVERITY_CONFIG[severity] ?? SEVERITY_CONFIG.low;
  return (
    <span style={{
      fontSize: 11, fontWeight: 600, padding: '2px 8px',
      borderRadius: 99, border: `1px solid ${cfg.border}`,
      background: cfg.bg, color: cfg.color,
      textTransform: 'uppercase', letterSpacing: '0.4px',
      whiteSpace: 'nowrap',
    }}>
      {cfg.label}
    </span>
  );
}

function formatDuration(startedAt) {
  const diff = Date.now() - new Date(startedAt).getTime();
  const m = Math.floor(diff / 60000);
  if (m < 1)  return '< 1m';
  if (m < 60) return `${m}m`;
  const h = Math.floor(m / 60);
  if (h < 24) return `${h}h ${m % 60}m`;
  return `${Math.floor(h / 24)}d ${h % 24}h`;
}

function formatDateTime(ts) {
  if (!ts) return '—';
  return new Date(ts).toLocaleString('id-ID', {
    day: '2-digit', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit',
  });
}

function formatDurationResolved(startedAt, resolvedAt) {
  if (!resolvedAt) return '—';
  const diff = new Date(resolvedAt) - new Date(startedAt);
  const m = Math.floor(diff / 60000);
  if (m < 60) return `${m}m`;
  const h = Math.floor(m / 60);
  return `${h}h ${m % 60}m`;
}

/* ─────────────────────────────────────
   Add Incident Modal
───────────────────────────────────── */
function AddIncidentModal({ onClose, onSubmit }) {
  const [form, setForm] = useState({
    title: '', description: '', severity: 'medium', created_by: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError]     = useState(null);

  function set(field, val) { setForm((p) => ({ ...p, [field]: val })); }

  async function handleSubmit() {
    if (!form.title.trim())      { setError('Judul wajib diisi.'); return; }
    if (!form.created_by.trim()) { setError('Nama pelapor wajib diisi.'); return; }
    setLoading(true);
    setError(null);
    try {
      await onSubmit({ ...form });
      onClose();
    } catch (err) {
      setError(err.message ?? 'Gagal membuat incident.');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="inc-overlay" onClick={() => !loading && onClose()}>
      <div className="inc-modal" onClick={(e) => e.stopPropagation()}>
        <div className="inc-modal__header">
          <span className="inc-modal__title">➕ Add Incident</span>
          <button className="inc-modal__close" onClick={onClose} disabled={loading}>×</button>
        </div>

        <div className="inc-modal__body">
          {error && <div className="inc-modal__error">{error}</div>}

          <label className="inc-field__label">Judul Gangguan *</label>
          <input
            className="inc-field__input"
            placeholder="e.g. WhatsApp gateway down"
            value={form.title}
            onChange={(e) => set('title', e.target.value)}
            disabled={loading}
          />

          <label className="inc-field__label">Deskripsi</label>
          <textarea
            className="inc-field__input inc-field__textarea"
            placeholder="Detail gangguan yang terjadi…"
            value={form.description}
            onChange={(e) => set('description', e.target.value)}
            disabled={loading}
            rows={3}
          />

          <label className="inc-field__label">Severity</label>
          <div className="inc-severity-grid">
            {Object.entries(SEVERITY_CONFIG).map(([key, cfg]) => (
              <button
                key={key}
                className={`inc-severity-btn ${form.severity === key ? 'inc-severity-btn--active' : ''}`}
                style={form.severity === key
                  ? { borderColor: cfg.border, background: cfg.bg, color: cfg.color }
                  : {}}
                onClick={() => set('severity', key)}
                disabled={loading}
              >
                {cfg.label}
              </button>
            ))}
          </div>

          <label className="inc-field__label">Nama Pelapor *</label>
          <input
            className="inc-field__input"
            placeholder="e.g. Admin Ops"
            value={form.created_by}
            onChange={(e) => set('created_by', e.target.value)}
            disabled={loading}
          />
        </div>

        <div className="inc-modal__footer">
          <button className="inc-modal__cancel" onClick={onClose} disabled={loading}>
            Batal
          </button>
          <button className="inc-modal__submit" onClick={handleSubmit} disabled={loading}>
            {loading
              ? <><span className="spin">↻</span> Menyimpan…</>
              : 'Buat Incident'}
          </button>
        </div>
      </div>
    </div>
  );
}

/* ─────────────────────────────────────
   History Modal
───────────────────────────────────── */
function HistoryModal({ onClose, history, loading, filter }) {
  return (
    <div className="inc-overlay" onClick={onClose}>
      <div className="inc-modal inc-modal--wide" onClick={(e) => e.stopPropagation()}>
        <div className="inc-modal__header">
          <span className="inc-modal__title">📋 Incident History</span>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <span style={{ fontSize: 12, color: 'var(--text-muted)' }}>
              {filter.from} s/d {filter.to}
            </span>
            <button className="inc-modal__close" onClick={onClose}>×</button>
          </div>
        </div>

        <div className="inc-modal__body">
          {loading && <div className="inc-empty">Memuat history…</div>}

          {!loading && history.length === 0 && (
            <div className="inc-empty">Tidak ada incident pada periode ini.</div>
          )}

          {!loading && history.length > 0 && (
            <table className="inc-table">
              <thead>
                <tr>
                  <th>Judul</th>
                  <th>Severity</th>
                  <th>Mulai</th>
                  <th>Selesai</th>
                  <th>Durasi</th>
                  <th>Dibuat oleh</th>
                </tr>
              </thead>
              <tbody>
                {history.map((inc) => (
                  <tr key={inc.id}>
                    <td>
                      <div className="inc-table__title">{inc.title}</div>
                      {inc.description && (
                        <div className="inc-table__desc">{inc.description}</div>
                      )}
                    </td>
                    <td><SeverityBadge severity={inc.severity} /></td>
                    <td className="inc-table__ts">{formatDateTime(inc.started_at)}</td>
                    <td className="inc-table__ts">{formatDateTime(inc.resolved_at)}</td>
                    <td className="inc-table__ts">{formatDurationResolved(inc.started_at, inc.resolved_at)}</td>
                    <td className="inc-table__ts">{inc.created_by ?? '—'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
}

/* ─────────────────────────────────────
   Main Section
───────────────────────────────────── */
export function IncidentSection({ filter }) {
  const {
    active, history, loading, histLoading,
    loadHistory, addIncident, resolve,
  } = useIncidents(filter);

  const [showAdd,     setShowAdd]     = useState(false);
  const [showHistory, setShowHistory] = useState(false);
  const [resolvingId, setResolvingId] = useState(null);
  const [toast, setToast]             = useState(null);

  const hasActive = active.length > 0;

  async function handleResolve(id) {
    setResolvingId(id);
    try {
      await resolve(id);
      setToast({ type: 'success', text: 'Incident berhasil di-resolve.' });
    } catch (err) {
      setToast({ type: 'error', text: err.message ?? 'Gagal resolve incident.' });
    } finally {
      setResolvingId(null);
    }
  }

  async function handleOpenHistory() {
    setShowHistory(true);
    await loadHistory();
  }

  async function handleAddSubmit(payload) {
    await addIncident(payload);
    setToast({ type: 'success', text: 'Incident berhasil dibuat.' });
  }

  return (
    <>
      <section className="inc-section">
        {/* ── Banner ── */}
        <div className={`inc-banner ${hasActive ? 'inc-banner--active' : 'inc-banner--clear'}`}>
          <div className="inc-banner__left">
            <span className={`inc-banner__dot ${hasActive ? 'inc-banner__dot--pulse' : ''}`} />
            <div>
              <div className="inc-banner__count">
                {loading ? '…' : `${active.length} Active`}
              </div>
              <div className="inc-banner__sub">
                {hasActive
                  ? 'Ada gangguan sistem yang sedang berlangsung'
                  : 'All systems operational. No active incidents.'}
              </div>
            </div>
          </div>
          <div className="inc-banner__actions">
            <button className="inc-btn inc-btn--ghost" onClick={handleOpenHistory}>
              📋 Show History
            </button>
            <button className="inc-btn inc-btn--add" onClick={() => setShowAdd(true)}>
              ＋ Add Incident
            </button>
          </div>
        </div>

        {/* ── Active incidents — tabel (selalu tampil) ── */}
        <div className="inc-active-table-wrap">
          <table className="inc-active-table">
            <thead>
              <tr>
                <th>Judul</th>
                <th>Severity</th>
                <th>Mulai</th>
                <th>Durasi</th>
                <th>Dibuat oleh</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={6} style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '20px', fontSize: 13 }}>
                    Memuat data…
                  </td>
                </tr>
              ) : active.length === 0 ? (
                <tr>
                  <td colSpan={6} style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '20px', fontSize: 13 }}>
                    — Tidak ada incident aktif —
                  </td>
                </tr>
              ) : active.map((inc) => {
                  const isResolving = resolvingId === inc.id;
                  return (
                    <tr key={inc.id}>
                      <td>
                        <div className="inc-active-table__title">{inc.title}</div>
                        {inc.description && (
                          <div className="inc-active-table__desc">{inc.description}</div>
                        )}
                      </td>
                      <td><SeverityBadge severity={inc.severity} /></td>
                      <td className="inc-active-table__ts">{formatDateTime(inc.started_at)}</td>
                      <td className="inc-active-table__ts">{formatDuration(inc.started_at)}</td>
                      <td className="inc-active-table__ts">{inc.created_by ?? '—'}</td>
                      <td>
                        <button
                          className="inc-resolve-btn"
                          onClick={() => handleResolve(inc.id)}
                          disabled={isResolving}
                          title="Tandai selesai"
                        >
                          {isResolving
                            ? <><span className="spin">↻</span> Resolving…</>
                            : <>✓ Resolve</>}
                        </button>
                      </td>
                    </tr>
                  );
              })}
            </tbody>
          </table>
        </div>
      </section>

      {/* ── Toast ── */}
      {toast && (
        <div
          className={`inc-toast inc-toast--${toast.type}`}
          onClick={() => setToast(null)}
        >
          {toast.type === 'success' ? '✓' : '✕'} {toast.text}
        </div>
      )}

      {/* ── Modals ── */}
      {showAdd && (
        <AddIncidentModal
          onClose={() => setShowAdd(false)}
          onSubmit={handleAddSubmit}
        />
      )}
      {showHistory && (
        <HistoryModal
          onClose={() => setShowHistory(false)}
          history={history}
          loading={histLoading}
          filter={filter}
        />
      )}

      <style>{`
        /* ── Section wrapper ── */
        .inc-section {
          margin: 0;
          display: flex;
          flex-direction: column;
          gap: var(--space-3);
        }

        /* ── Banner ── */
        .inc-banner {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 14px 20px;
          border-radius: var(--radius-md);
          border: 1px solid;
          gap: 16px;
        }
        .inc-banner--clear {
          background: rgba(34,197,94,0.06);
          border-color: rgba(34,197,94,0.2);
        }
        .inc-banner--active {
          background: rgba(239,68,68,0.07);
          border-color: rgba(239,68,68,0.3);
          animation: incPulseBorder 2s ease-in-out infinite;
        }
        @keyframes incPulseBorder {
          0%,100% { box-shadow: 0 0 0 0 rgba(239,68,68,0.2); }
          50%      { box-shadow: 0 0 0 6px rgba(239,68,68,0); }
        }
        .inc-banner__left { display: flex; align-items: center; gap: 12px; }
        .inc-banner__dot {
          width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0;
          background: var(--color-success);
        }
        .inc-banner--active .inc-banner__dot { background: var(--color-danger); }
        .inc-banner__dot--pulse { animation: dotPulse 1.2s ease-in-out infinite; }
        @keyframes dotPulse {
          0%,100% { opacity:1; transform:scale(1); box-shadow:0 0 0 0 rgba(239,68,68,0.5); }
          50%      { opacity:0.8; transform:scale(1.2); box-shadow:0 0 0 5px rgba(239,68,68,0); }
        }
        .inc-banner__count {
          font-size: 15px; font-weight: 700; color: var(--text-primary);
        }
        .inc-banner--active .inc-banner__count { color: var(--color-danger); }
        .inc-banner__sub { font-size: 12px; color: var(--text-muted); margin-top: 2px; }
        .inc-banner__actions { display: flex; gap: 8px; flex-shrink: 0; }

        /* ── Banner buttons ── */
        .inc-btn {
          display: flex; align-items: center; gap: 6px;
          padding: 6px 14px; font-size: 13px; font-weight: 500;
          border-radius: var(--radius-sm); cursor: pointer;
          transition: opacity 0.15s; border: 1px solid;
        }
        .inc-btn:hover { opacity: 0.8; }
        .inc-btn--ghost {
          background: transparent;
          border-color: var(--border-default);
          color: var(--text-secondary);
        }
        .inc-btn--add {
          background: rgba(239,68,68,0.1);
          border-color: rgba(239,68,68,0.3);
          color: var(--color-danger);
        }

        /* ── Active incidents table ── */
        .inc-active-table-wrap {
          background: var(--bg-card);
          border: 1px solid rgba(239,68,68,0.2);
          border-radius: var(--radius-md);
          overflow: hidden;
        }
        .inc-active-table {
          width: 100%; border-collapse: collapse; font-size: 13px;
        }
        .inc-active-table th {
          text-align: left; font-size: 11px; font-weight: 500;
          color: var(--text-muted); text-transform: uppercase;
          letter-spacing: 0.4px; padding: 10px 14px;
          background: rgba(239,68,68,0.04);
          border-bottom: 1px solid rgba(239,68,68,0.12);
        }
        .inc-active-table td {
          padding: 11px 14px;
          border-bottom: 1px solid var(--border-subtle);
          vertical-align: middle; color: var(--text-secondary);
        }
        .inc-active-table tr:last-child td { border-bottom: none; }
        .inc-active-table tr:hover td { background: rgba(255,255,255,0.02); }
        .inc-active-table__title { font-weight: 500; color: var(--text-primary); font-size: 13px; }
        .inc-active-table__desc { font-size: 11px; color: var(--text-muted); margin-top: 2px; }
        .inc-active-table__ts { font-size: 12px; white-space: nowrap; }

        /* ── Resolve button ── */
        .inc-resolve-btn {
          display: inline-flex; align-items: center; gap: 5px;
          padding: 5px 14px; font-size: 12px; font-weight: 600;
          background: rgba(34,197,94,0.12);
          border: 1px solid rgba(34,197,94,0.4);
          border-radius: var(--radius-sm);
          color: var(--color-success); cursor: pointer;
          transition: background 0.15s, box-shadow 0.15s, transform 0.1s;
          white-space: nowrap; user-select: none;
        }
        .inc-resolve-btn:hover:not(:disabled) {
          background: rgba(34,197,94,0.22);
          box-shadow: 0 0 0 3px rgba(34,197,94,0.15);
          transform: translateY(-1px);
        }
        .inc-resolve-btn:active:not(:disabled) { transform: translateY(0); box-shadow: none; }
        .inc-resolve-btn:disabled { opacity: 0.5; cursor: not-allowed; }

        /* ── Modal overlay & container ── */
        .inc-overlay {
          position: fixed; inset: 0; z-index: 1000;
          background: rgba(0,0,0,0.6);
          display: flex; align-items: center; justify-content: center;
          animation: fadeIn 0.15s ease;
        }
        .inc-modal {
          background: var(--bg-card);
          border: 1px solid var(--border-default);
          border-radius: var(--radius-lg);
          width: 460px; max-width: calc(100vw - 32px); max-height: 85vh;
          display: flex; flex-direction: column;
          box-shadow: var(--shadow-elevated);
          animation: slideUp 0.18s ease;
          overflow: hidden;
        }
        .inc-modal--wide { width: 800px; }
        .inc-modal__header {
          display: flex; align-items: center; justify-content: space-between;
          padding: 16px 20px 12px;
          border-bottom: 1px solid var(--border-subtle);
          flex-shrink: 0;
        }
        .inc-modal__title { font-size: 15px; font-weight: 600; color: var(--text-primary); }
        .inc-modal__close {
          background: none; border: none;
          color: var(--text-muted); font-size: 20px; cursor: pointer; line-height: 1;
        }
        .inc-modal__close:hover { color: var(--text-primary); }
        .inc-modal__body {
          padding: 16px 20px; display: flex; flex-direction: column;
          gap: 10px; overflow-y: auto; flex: 1;
        }
        .inc-modal__error {
          padding: 8px 12px;
          background: var(--color-danger-bg);
          border: 1px solid var(--color-danger-border);
          border-radius: var(--radius-sm);
          color: var(--color-danger); font-size: 13px;
        }
        .inc-modal__footer {
          display: flex; justify-content: flex-end; gap: 8px;
          padding: 12px 20px 16px;
          border-top: 1px solid var(--border-subtle);
          flex-shrink: 0;
        }
        .inc-modal__cancel {
          padding: 7px 16px; font-size: 13px; background: none;
          border: 1px solid var(--border-default);
          border-radius: var(--radius-sm);
          color: var(--text-secondary); cursor: pointer;
        }
        .inc-modal__cancel:hover:not(:disabled) {
          border-color: var(--border-strong); color: var(--text-primary);
        }
        .inc-modal__submit {
          display: flex; align-items: center; gap: 6px;
          padding: 7px 16px; font-size: 13px; font-weight: 500;
          background: rgba(239,68,68,0.1);
          border: 1px solid rgba(239,68,68,0.3);
          border-radius: var(--radius-sm);
          color: var(--color-danger); cursor: pointer;
        }
        .inc-modal__submit:hover:not(:disabled) { opacity: 0.8; }
        .inc-modal__submit:disabled { opacity: 0.6; cursor: not-allowed; }

        /* ── Form fields ── */
        .inc-field__label {
          font-size: 12px; font-weight: 500; color: var(--text-muted);
          text-transform: uppercase; letter-spacing: 0.4px; margin-bottom: -4px;
        }
        .inc-field__input {
          width: 100%; padding: 8px 12px; font-size: 13px;
          background: var(--bg-surface);
          border: 1px solid var(--border-default);
          border-radius: var(--radius-sm);
          color: var(--text-primary); font-family: var(--font-sans);
          transition: border-color 0.15s; box-sizing: border-box;
        }
        .inc-field__input:focus { outline: none; border-color: var(--accent-primary); }
        .inc-field__input:disabled { opacity: 0.5; }
        .inc-field__textarea { resize: vertical; min-height: 72px; }

        /* ── Severity picker ── */
        .inc-severity-grid {
          display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px;
        }
        .inc-severity-btn {
          padding: 7px 0; font-size: 12px; font-weight: 500;
          background: var(--bg-surface);
          border: 1px solid var(--border-default);
          border-radius: var(--radius-sm);
          color: var(--text-secondary); cursor: pointer; transition: all 0.15s;
        }
        .inc-severity-btn:hover:not(:disabled) {
          border-color: var(--border-strong); color: var(--text-primary);
        }
        .inc-severity-btn:disabled { opacity: 0.5; cursor: not-allowed; }

        /* ── History table ── */
        .inc-table {
          width: 100%; border-collapse: collapse; font-size: 13px;
        }
        .inc-table th {
          text-align: left; font-size: 11px; font-weight: 500;
          color: var(--text-muted); text-transform: uppercase;
          letter-spacing: 0.4px; padding: 6px 10px;
          border-bottom: 1px solid var(--border-subtle);
        }
        .inc-table td {
          padding: 10px; border-bottom: 1px solid var(--border-subtle);
          vertical-align: top; color: var(--text-secondary);
        }
        .inc-table tr:last-child td { border-bottom: none; }
        .inc-table tr:hover td { background: rgba(255,255,255,0.02); }
        .inc-table__title { font-weight: 500; color: var(--text-primary); }
        .inc-table__desc { font-size: 11px; color: var(--text-muted); margin-top: 2px; }
        .inc-table__ts { font-size: 12px; white-space: nowrap; }

        /* ── Toast ── */
        .inc-toast {
          position: fixed; bottom: 24px; right: 24px; z-index: 9999;
          padding: 10px 16px; border-radius: var(--radius-md);
          font-size: 13px; font-weight: 500; cursor: pointer;
          animation: slideIn 0.2s ease; box-shadow: var(--shadow-elevated);
        }
        .inc-toast--success {
          background: var(--color-success-bg);
          border: 1px solid var(--color-success-border);
          color: var(--color-success);
        }
        .inc-toast--error {
          background: var(--color-danger-bg);
          border: 1px solid var(--color-danger-border);
          color: var(--color-danger);
        }

        /* ── Empty state ── */
        .inc-empty {
          text-align: center; padding: 32px;
          color: var(--text-muted); font-size: 13px;
        }

        /* ── Animations ── */
        .spin { display: inline-block; animation: spin 0.8s linear infinite; }
        @keyframes spin { from{transform:rotate(0deg)} to{transform:rotate(360deg)} }
        @keyframes fadeIn { from{opacity:0} to{opacity:1} }
        @keyframes slideUp {
          from{opacity:0;transform:translateY(10px)} to{opacity:1;transform:translateY(0)}
        }
        @keyframes slideIn {
          from{opacity:0;transform:translateY(12px)} to{opacity:1;transform:translateY(0)}
        }
      `}</style>
    </>
  );
}