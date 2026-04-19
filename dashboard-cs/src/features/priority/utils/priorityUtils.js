export const PRIORITY_CONFIG = {
  roaming: {
    label: 'Roaming',
    color: '#f97316',
    bg: 'rgba(249,115,22,0.12)',
    border: 'rgba(249,115,22,0.25)',
    icon: '📡',
  },
  extra_quota: {
    label: 'Extra Quota',
    color: '#a78bfa',
    bg: 'rgba(167,139,250,0.12)',
    border: 'rgba(167,139,250,0.25)',
    icon: '📶',
  },
  cc: {
    label: 'CC',
    color: '#60a5fa',
    bg: 'rgba(96,165,250,0.12)',
    border: 'rgba(96,165,250,0.25)',
    icon: '🎧',
  },
  vip: {
    label: 'VIP',
    color: '#facc15',
    bg: 'rgba(250,204,21,0.12)',
    border: 'rgba(250,204,21,0.25)',
    icon: '👑',
  },
  p1: {
    label: 'P1',
    color: '#ef4444',
    bg: 'rgba(239,68,68,0.12)',
    border: 'rgba(239,68,68,0.3)',
    icon: '🔴',
  },
  urgent: {
    label: 'Urgent',
    color: '#f43f5e',
    bg: 'rgba(244,63,94,0.12)',
    border: 'rgba(244,63,94,0.3)',
    icon: '⚡',
  },
};

export function getUrgencyLevel(key) {
  if (['p1', 'urgent'].includes(key)) return 'critical';
  if (['roaming', 'vip'].includes(key)) return 'high';
  return 'normal';
}
