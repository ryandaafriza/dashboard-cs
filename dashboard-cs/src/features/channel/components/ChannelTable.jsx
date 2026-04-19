import React from 'react';

export function ChannelTable({ title, rows, color }) {
  const isEmpty = !rows || rows.length === 0;

  return (
    <div className="channel-table">
      <div
        className="channel-table__title"
        style={{ borderLeft: `2px solid ${color}` }}
      >
        {title}
      </div>

      {isEmpty ? (
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
            {rows.map((row, i) => (
              <tr key={i} className="channel-table__row">
                <td className="channel-table__td channel-table__td--name" title={row.name}>
                  {row.name}
                </td>
                <td className="channel-table__td channel-table__td--num">{row.interactions}</td>
                <td className="channel-table__td channel-table__td--num">{row.tickets}</td>
                <td className="channel-table__td channel-table__td--num">
                  <span
                    className="channel-table__fcr"
                    style={{
                      color: row.fcr >= 80 ? '#22c55e' : row.fcr >= 50 ? '#f59e0b' : 'var(--text-muted)',
                    }}
                  >
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
