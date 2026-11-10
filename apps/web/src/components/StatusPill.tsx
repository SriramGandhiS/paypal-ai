import React from 'react';

interface StatusPillProps {
  status: string;
  type?: 'provider' | 'app' | 'case' | 'tier';
}

export const StatusPill: React.FC<StatusPillProps> = ({ status, type = 'case' }) => {
  let bg = 'bg-slate-800 text-slate-300 border-slate-700';
  let dot = 'bg-slate-400';

  const s = status.toUpperCase();

  if (s === 'CAPTURE_COMPLETED' || s === 'SYNCHRONISED' || s === 'RESOLVED' || s === 'CREDITED' || s === 'SAFE') {
    bg = 'bg-emerald-950/60 text-emerald-300 border-emerald-800/80';
    dot = 'bg-emerald-400';
  } else if (s === 'API_TIMEOUT' || s === 'RECONCILING' || s === 'OBSERVING' || s === 'CREDIT_PENDING' || s === 'LOW') {
    bg = 'bg-amber-950/60 text-amber-300 border-amber-800/80';
    dot = 'bg-amber-400';
  } else if (s === 'INCONSISTENT' || s === 'NEEDS_HUMAN_REVIEW' || s === 'CAPTURE_DENIED' || s === 'HIGH') {
    bg = 'bg-rose-950/60 text-rose-300 border-rose-800/80';
    dot = 'bg-rose-400';
  } else if (s === 'UNKNOWN' || s === 'NOT_OBSERVABLE') {
    // Invariant: UNKNOWN is never red, never failed
    bg = 'bg-slate-900/80 text-slate-300 border-slate-700';
    dot = 'bg-slate-400';
  } else if (s.startsWith('TIER 1')) {
    bg = 'bg-indigo-950/80 text-indigo-300 border-indigo-700';
    dot = 'bg-indigo-400';
  }

  return (
    <span className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium border ${bg}`}>
      <span className={`w-1.5 h-1.5 rounded-full ${dot}`}></span>
      <span>{status.replace(/_/g, ' ')}</span>
    </span>
  );
};
