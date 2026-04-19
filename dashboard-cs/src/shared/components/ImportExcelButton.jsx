import React, { useRef, useState } from 'react';
import { importExcel } from '../../services/importService';

export function ImportExcelButton({ onImportSuccess }) {
  const inputRef = useRef(null);
  const [loading, setLoading] = useState(false);

  function handleClick() {
    inputRef.current?.click();
  }

  async function handleFileChange(e) {
    const file = e.target.files?.[0];
    if (!file) return;

    // Reset input agar file yang sama bisa diupload ulang
    e.target.value = '';

    setLoading(true);
    try {
      await importExcel(file);
      onImportSuccess({ type: 'success', message: `File "${file.name}" berhasil diupload.` });
    } catch (err) {
      onImportSuccess({ type: 'error', message: err.message ?? 'Gagal mengupload file.' });
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      <input
        ref={inputRef}
        type="file"
        accept=".xlsx,.xls"
        style={{ display: 'none' }}
        onChange={handleFileChange}
      />
      <button
        className="topbar-btn topbar-btn--import"
        onClick={handleClick}
        disabled={loading}
        title="Import Excel"
      >
        {loading ? (
          <span className="spin">↻</span>
        ) : (
          <span>⬆</span>
        )}
        {loading ? 'Uploading…' : 'Import Excel'}
      </button>
    </>
  );
}