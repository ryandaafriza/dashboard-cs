import React, { useState } from 'react';

export function SyncButton({ onSync, needsSync }) {
  const [loading, setLoading] = useState(false);

  async function handleSync() {
    setLoading(true);
    try {
      await onSync();
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      <button
        className={`topbar-btn sync-state-btn ${needsSync ? 'sync-state-btn--warning' : 'sync-state-btn--synced'}`}
        onClick={handleSync}
        disabled={loading}
        title={needsSync ? 'Data belum ter-sync, klik untuk sync' : 'Klik untuk sync ulang'}
      >
        <span className={loading ? 'spin' : needsSync ? 'pulse-dot' : ''}>
          {loading ? '↻' : needsSync ? '' : '↻'}
        </span>
        {loading ? 'Syncing…' : needsSync ? 'Data belum ter-sync' : 'Sync'}
      </button>

      <style>{`
        .sync-state-btn {
          transition: background 0.3s, border-color 0.3s, color 0.3s;
        }
        .sync-state-btn--warning {
          background: var(--color-warning-bg, rgba(245,158,11,0.1));
          border-color: var(--color-warning-border, rgba(245,158,11,0.3)) !important;
          color: var(--color-warning, #f59e0b);
          animation: pulse-border-warn 1.5s ease-in-out infinite;
        }
        .sync-state-btn--synced {
          background: var(--color-info-bg, rgba(59,130,246,0.1));
          border-color: var(--color-info-border, rgba(59,130,246,0.25)) !important;
          color: var(--color-info, #3b82f6);
        }
        .sync-state-btn--synced:hover:not(:disabled) { opacity: 0.8; }
        @keyframes pulse-border-warn {
          0%,100% { box-shadow: 0 0 0 0 rgba(245,158,11,0.3); }
          50%      { box-shadow: 0 0 0 4px rgba(245,158,11,0); }
        }
        .pulse-dot {
          display: inline-block;
          width: 7px; height: 7px;
          border-radius: 50%;
          background: var(--color-warning, #f59e0b);
          animation: pulseDot 1.2s ease-in-out infinite;
        }
        .spin { display: inline-block; animation: spin 0.8s linear infinite; }
        @keyframes pulseDot {
          0%,100% { opacity:1; transform:scale(1); }
          50%      { opacity:0.4; transform:scale(0.7); }
        }
        @keyframes spin { from{transform:rotate(0deg)} to{transform:rotate(360deg)} }
      `}</style>
    </>
  );
}