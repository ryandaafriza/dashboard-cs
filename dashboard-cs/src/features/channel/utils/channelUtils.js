export const CHANNEL_CONFIG = {
  email: {
    label: 'Email',
    icon: '✉',
    color: '#60a5fa',
    gradient: 'linear-gradient(135deg, rgba(96,165,250,0.08) 0%, transparent 60%)',
  },
  whatsapp: {
    label: 'WhatsApp',
    icon: '💬',
    color: '#22c55e',
    gradient: 'linear-gradient(135deg, rgba(34,197,94,0.08) 0%, transparent 60%)',
  },
  social_media: {
    label: 'Social Media',
    icon: '📱',
    color: '#f472b6',
    gradient: 'linear-gradient(135deg, rgba(244,114,182,0.08) 0%, transparent 60%)',
  },
  live_chat: {
    label: 'Live Chat',
    icon: '💡',
    color: '#fb923c',
    gradient: 'linear-gradient(135deg, rgba(251,146,60,0.08) 0%, transparent 60%)',
  },
  call_center: {
    label: 'Call Center 188',
    icon: '📞',
    color: '#a78bfa',
    gradient: 'linear-gradient(135deg, rgba(167,139,250,0.08) 0%, transparent 60%)',
  },
};

export function getSLAStatus(sla) {
  if (sla >= 80) return { level: 'good',    color: '#22c55e', bg: 'rgba(34,197,94,0.12)',   border: 'rgba(34,197,94,0.25)' };
  if (sla >= 60) return { level: 'medium',  color: '#f59e0b', bg: 'rgba(245,158,11,0.12)',  border: 'rgba(245,158,11,0.25)' };
  if (sla > 0)   return { level: 'low',     color: '#ef4444', bg: 'rgba(239,68,68,0.12)',   border: 'rgba(239,68,68,0.25)' };
  return              { level: 'critical', color: '#ef4444', bg: 'rgba(239,68,68,0.08)',   border: 'rgba(239,68,68,0.2)' };
}
