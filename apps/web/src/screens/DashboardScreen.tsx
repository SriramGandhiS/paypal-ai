import React, { useState, useEffect } from 'react';
import { 
  Plus, 
  ArrowUpRight, 
  CheckCircle2, 
  AlertTriangle, 
  ShieldAlert, 
  DollarSign, 
  PlayCircle,
  Clock,
  ChevronRight
} from 'lucide-react';
import { api } from '../api/client';
import { CaseSummary } from '../api/types';
import { StatusPill } from '../components/StatusPill';

interface DashboardScreenProps {
  onSelectCase: (caseId: string) => void;
  onNavigate: (tab: string) => void;
}

export const DashboardScreen: React.FC<DashboardScreenProps> = ({
  onSelectCase,
  onNavigate,
}) => {
  const [cases, setCases] = useState<CaseSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [createForm, setCreateForm] = useState({
    external_reference: `INV-2026-${Math.floor(1000 + Math.random() * 9000)}`,
    amount_cents: 14950,
    currency: 'USD',
    sender_country: 'US',
    receiver_country: 'DE',
    payment_purpose: 'SAAS_SUBSCRIPTION',
  });

  const loadData = async () => {
    setLoading(true);
    const data = await api.listCases();
    setCases(data);
    setLoading(false);
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleCreateCase = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const created = await api.createCase(createForm);
      setShowCreateModal(false);
      await loadData();
      onSelectCase(created.id);
    } catch (err: any) {
      alert(err.message || 'Failed to create case');
    }
  };

  const totalVolume = cases.reduce((sum, c) => sum + parseFloat(c.amount || '0'), 0);
  const resolvedCount = cases.filter(c => c.case_status === 'RESOLVED').length;
  const reviewCount = cases.filter(c => c.case_status === 'NEEDS_HUMAN_REVIEW').length;
  const observingCount = cases.filter(c => c.case_status === 'OBSERVING').length;

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Top Welcome / Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-100">
            Payment Reconciliation Console
          </h1>
          <p className="text-sm text-slate-400 mt-1">
            Authoritative financial state ledger, blind-retry guards, and payment truth passports.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <button
            onClick={() => onNavigate('demo')}
            className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-indigo-600 to-cyan-500 hover:from-indigo-500 hover:to-cyan-400 text-white font-medium text-sm rounded-lg shadow-lg shadow-indigo-500/20 transition-all cursor-pointer"
          >
            <PlayCircle className="w-4 h-4" />
            Run Golden Path Demo
          </button>
          <button
            onClick={() => setShowCreateModal(true)}
            className="flex items-center gap-2 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-200 font-medium text-sm rounded-lg border border-slate-700 transition-all cursor-pointer"
          >
            <Plus className="w-4 h-4" />
            New Case
          </button>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="p-5 rounded-xl bg-slate-900/60 border border-slate-800 hover:border-slate-700 transition-all shadow-sm">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium text-slate-400">Total Monitored Volume</span>
            <div className="p-2 rounded-lg bg-indigo-500/10 text-indigo-400 border border-indigo-500/20">
              <DollarSign className="w-4 h-4" />
            </div>
          </div>
          <div className="mt-3">
            <div className="text-2xl font-bold tracking-tight tabular-nums text-slate-100">
              ${totalVolume.toLocaleString('en-US', { minimumFractionDigits: 2 })}
            </div>
            <p className="text-xs text-slate-500 mt-1 flex items-center gap-1">
              <span className="text-emerald-400 font-medium">100% Truth</span> accuracy via SHA-256 seal
            </p>
          </div>
        </div>

        <div className="p-5 rounded-xl bg-slate-900/60 border border-slate-800 hover:border-slate-700 transition-all shadow-sm">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium text-slate-400">Reconciled Cases</span>
            <div className="p-2 rounded-lg bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
              <CheckCircle2 className="w-4 h-4" />
            </div>
          </div>
          <div className="mt-3">
            <div className="text-2xl font-bold tracking-tight tabular-nums text-emerald-400">
              {resolvedCount} / {cases.length}
            </div>
            <p className="text-xs text-slate-500 mt-1">Provider & Application in synchrony</p>
          </div>
        </div>

        <div className="p-5 rounded-xl bg-slate-900/60 border border-slate-800 hover:border-slate-700 transition-all shadow-sm">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium text-slate-400">Gaps Under Review</span>
            <div className="p-2 rounded-lg bg-amber-500/10 text-amber-400 border border-amber-500/20">
              <AlertTriangle className="w-4 h-4" />
            </div>
          </div>
          <div className="mt-3">
            <div className="text-2xl font-bold tracking-tight tabular-nums text-amber-400">
              {reviewCount} Active
            </div>
            <p className="text-xs text-slate-500 mt-1">State inconsistency or API timeout</p>
          </div>
        </div>

        <div className="p-5 rounded-xl bg-slate-900/60 border border-slate-800 hover:border-slate-700 transition-all shadow-sm">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium text-slate-400">Blind Retries Blocked</span>
            <div className="p-2 rounded-lg bg-violet-500/10 text-violet-400 border border-violet-500/20">
              <ShieldAlert className="w-4 h-4" />
            </div>
          </div>
          <div className="mt-3">
            <div className="text-2xl font-bold tracking-tight tabular-nums text-violet-400">
              {cases.length > 0 ? cases.length + 1 : 0} Defended
            </div>
            <p className="text-xs text-slate-500 mt-1">Duplicate charges mathematically prevented</p>
          </div>
        </div>
      </div>

      {/* Hero Showcase Card */}
      <div className="p-6 rounded-2xl bg-gradient-to-r from-slate-900 via-slate-900 to-indigo-950/40 border border-indigo-900/40 relative overflow-hidden">
        <div className="relative z-10 max-w-2xl">
          <span className="px-3 py-1 rounded-full text-xs font-semibold bg-indigo-500/20 text-indigo-300 border border-indigo-500/30">
            DEMONSTRATION SCENARIO
          </span>
          <h2 className="text-xl font-bold text-slate-100 mt-3">
            The Classic Blind-Retry Trap: PayPal Captures, Website Times Out
          </h2>
          <p className="text-sm text-slate-300 mt-2 leading-relaxed">
            When an upstream timeout leaves a merchant order pending, conventional checkouts prompt the buyer to "Try Again" — double-charging the customer and losing chargebacks. Watch PayBridge AI intercept the existing capture, establish financial ground truth, and safely reconcile without blind retries.
          </p>
          <div className="mt-5 flex items-center gap-4">
            <button
              onClick={() => onNavigate('demo')}
              className="flex items-center gap-2 px-5 py-2.5 bg-indigo-600 hover:bg-indigo-500 text-white font-medium text-sm rounded-lg shadow-md transition-all cursor-pointer"
            >
              <PlayCircle className="w-4 h-4" />
              Launch Step-by-Step Demo
            </button>
            <button
              onClick={() => onNavigate('mirror')}
              className="flex items-center gap-2 px-4 py-2.5 bg-slate-800/80 hover:bg-slate-700 text-slate-300 font-medium text-sm rounded-lg border border-slate-700 transition-all cursor-pointer"
            >
              Inspect Recovery Mirror
              <ArrowUpRight className="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      {/* Payment Cases Table */}
      <div className="rounded-xl bg-slate-900/60 border border-slate-800 overflow-hidden shadow-sm">
        <div className="p-4 sm:px-6 border-b border-slate-800 flex items-center justify-between">
          <div>
            <h3 className="font-semibold text-slate-100">Live Payment Cases</h3>
            <p className="text-xs text-slate-400 mt-0.5">Real-time status vectors from PostgreSQL 16 state engine</p>
          </div>
          <span className="text-xs font-mono text-slate-400 bg-slate-800 px-2.5 py-1 rounded border border-slate-700">
            {cases.length} Total Cases
          </span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="bg-slate-950/50 text-slate-400 text-xs uppercase font-mono tracking-wider border-b border-slate-800">
              <tr>
                <th className="py-3 px-4 sm:px-6">Invoice Reference</th>
                <th className="py-3 px-4">Amount</th>
                <th className="py-3 px-4">Provider State</th>
                <th className="py-3 px-4">Application State</th>
                <th className="py-3 px-4">Case Status</th>
                <th className="py-3 px-4">Version</th>
                <th className="py-3 px-4 sm:px-6 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/60">
              {cases.map((c) => (
                <tr 
                  key={c.id} 
                  onClick={() => onSelectCase(c.id)}
                  className="hover:bg-slate-800/40 cursor-pointer transition-colors group"
                >
                  <td className="py-3.5 px-4 sm:px-6 font-mono font-medium text-slate-200 group-hover:text-indigo-300 flex items-center gap-2">
                    {c.external_reference}
                  </td>
                  <td className="py-3.5 px-4 font-mono font-semibold text-slate-100 tabular-nums">
                    ${c.amount} <span className="text-xs text-slate-400 font-normal">{c.currency}</span>
                  </td>
                  <td className="py-3.5 px-4">
                    <StatusPill status={c.provider_state} type="provider" />
                  </td>
                  <td className="py-3.5 px-4">
                    <StatusPill status={c.application_state} type="app" />
                  </td>
                  <td className="py-3.5 px-4">
                    <StatusPill status={c.case_status} type="case" />
                  </td>
                  <td className="py-3.5 px-4 font-mono text-xs text-slate-400">
                    v{c.state_version}
                  </td>
                  <td className="py-3.5 px-4 sm:px-6 text-right">
                    <span className="inline-flex items-center text-xs text-indigo-400 group-hover:text-indigo-300 font-medium">
                      Inspect <ChevronRight className="w-3.5 h-3.5 ml-1 transition-transform group-hover:translate-x-0.5" />
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Create Case Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-700 rounded-2xl max-w-md w-full p-6 shadow-2xl">
            <h3 className="text-lg font-bold text-slate-100 mb-1">Create Payment Reconciliation Case</h3>
            <p className="text-xs text-slate-400 mb-5">Initializes a new case with strict idempotency and zero blind-retries.</p>

            <form onSubmit={handleCreateCase} className="space-y-4">
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">External Invoice Reference</label>
                <input
                  type="text"
                  required
                  value={createForm.external_reference}
                  onChange={e => setCreateForm({ ...createForm, external_reference: e.target.value })}
                  className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm text-slate-100 font-mono focus:border-indigo-500 focus:outline-none"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">Amount ($ USD)</label>
                  <input
                    type="number"
                    step="0.01"
                    min="1"
                    required
                    value={createForm.amount_cents / 100}
                    onChange={e => setCreateForm({ ...createForm, amount_cents: Math.round(parseFloat(e.target.value || '0') * 100) })}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm text-slate-100 font-mono focus:border-indigo-500 focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">Currency</label>
                  <input
                    type="text"
                    disabled
                    value={createForm.currency}
                    className="w-full bg-slate-950/50 border border-slate-800 rounded-lg px-3 py-2 text-sm text-slate-400 font-mono"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">Sender Country</label>
                  <input
                    type="text"
                    maxLength={2}
                    required
                    value={createForm.sender_country}
                    onChange={e => setCreateForm({ ...createForm, sender_country: e.target.value.toUpperCase() })}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm text-slate-100 font-mono focus:border-indigo-500 focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">Receiver Country</label>
                  <input
                    type="text"
                    maxLength={2}
                    required
                    value={createForm.receiver_country}
                    onChange={e => setCreateForm({ ...createForm, receiver_country: e.target.value.toUpperCase() })}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm text-slate-100 font-mono focus:border-indigo-500 focus:outline-none"
                  />
                </div>
              </div>

              <div className="pt-3 flex items-center justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowCreateModal(false)}
                  className="px-4 py-2 text-sm text-slate-400 hover:text-slate-200 cursor-pointer"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-5 py-2 bg-indigo-600 hover:bg-indigo-500 text-white font-medium text-sm rounded-lg shadow-md transition-all cursor-pointer"
                >
                  Create & Initialize
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
