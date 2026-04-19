import React, { useEffect } from 'react';

export function Toast({ toast, onClose }) {
  useEffect(() => {
    if (!toast) return;
    const t = setTimeout(onClose, 4000);
    return () => clearTimeout(t);
  }, [toast, onClose]);

  if (!toast) return null;

  const isSuccess = toast.type === 'success';

  return (
    <div className={`toast toast--${toast.type}`} role="alert">
      <span className="toast__icon">{isSuccess ? '✓' : '✕'}</span>
      <span className="toast__message">{toast.message}</span>
      <button className="toast__close" onClick={onClose}>×</button>

      <style>{`
        .toast {
          position: fixed;
          bottom: 24px;
          right: 24px;
          z-index: 9999;
          display: flex;
          align-items: center;
          gap: 10px;
          padding: 12px 16px;
          border-radius: var(--radius-md, 10px);
          font-size: 14px;
          font-weight: 500;
          box-shadow: var(--shadow-elevated, 0 4px 24px rgba(0,0,0,0.5));
          animation: slideIn 0.2s ease;
          min-width: 280px;
          max-width: 420px;
        }
        .toast--success {
          background: var(--color-success-bg, rgba(34,197,94,0.12));
          border: 1px solid var(--color-success-border, rgba(34,197,94,0.3));
          color: var(--color-success, #22c55e);
        }
        .toast--error {
          background: var(--color-danger-bg, rgba(239,68,68,0.12));
          border: 1px solid var(--color-danger-border, rgba(239,68,68,0.3));
          color: var(--color-danger, #ef4444);
        }
        .toast__icon { font-size: 16px; flex-shrink: 0; }
        .toast__message { flex: 1; color: var(--text-primary, #e8edf5); }
        .toast__close {
          background: none; border: none; cursor: pointer;
          color: var(--text-muted, #4d5a70); font-size: 18px;
          line-height: 1; padding: 0 2px;
          flex-shrink: 0;
        }
        .toast__close:hover { color: var(--text-primary, #e8edf5); }
        @keyframes slideIn {
          from { opacity: 0; transform: translateY(12px); }
          to   { opacity: 1; transform: translateY(0); }
        }
      `}</style>
    </div>
  );
}