import React, { useState } from 'react';
import { 
  GitFork, 
  ShieldCheck, 
  AlertTriangle, 
  CheckCircle2, 
  XCircle, 
  DollarSign, 
  Clock, 
  ArrowRight, 
  TrendingDown, 
  Scale, 
  ShieldAlert,
  HelpCircle,
  Zap,
  RotateCcw
} from 'lucide-react';

interface Scenario {
  id: string;
  title: string;
  amount: string;
  currency: string;
  errorType: string;
  merchantRef: string;
}

const SCENARIOS: Scenario[] = [
  {
    id: 'sc-1',
    title: 'PayPal Capture Timeout on E-Commerce Order',
    amount: '149.50',
    currency: 'USD',
    errorType: '504 Gateway Timeout during capture execution',
    merchantRef: 'INV-2026-8941',
  },
  {
    id: 'sc-2',
    title: 'SaaS Subscription Upgrade Network Drop',
    amount: '580.00',
    currency: 'USD',
    errorType: 'Connection reset before webhook delivery',
    merchantRef: 'INV-2026-9022',
  },
  {
    id: 'sc-3',
    title: 'High-Value Enterprise Cloud Invoice',
    amount: '1,200.00',
    currency: 'USD',
    errorType: 'Asynchronous capture delayed by provider fraud check',
    merchantRef: 'INV-2026-9105',
  },
];

export const RecoveryMirrorScreen: React.FC = () => {
  const [selectedScenario, setSelectedScenario] = useState<Scenario>(SCENARIOS[0]);

  return (
    <div className="space-y-8 animate-fade-in max-w-6xl mx-auto pb-12">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <GitFork className="w-6 h-6 text-indigo-400" />
            <h1 className="text-2xl font-bold tracking-tight text-slate-100">
              Counterfactual Recovery Mirror
            </h1>
          </div>
          <p className="text-sm text-slate-400 mt-1">
            Compare the operational, financial, and customer outcomes of PayBridge automated truth resolution vs blind-retry defaults.
          </p>
        </div>

        {/* Scenario Selector */}
        <div className="flex items-center gap-2 bg-slate-900 border border-slate-800 rounded-xl p-1.5">
          <span className="text-xs font-medium text-slate-400 px-2">Scenario:</span>
          <select 
            value={selectedScenario.id}
            onChange={(e) => {
              const sc = SCENARIOS.find(s => s.id === e.target.value);
              if (sc) setSelectedScenario(sc);
            }}
            aria-label="Select Scenario"
            className="bg-slate-950 border border-slate-700/60 text-slate-200 text-xs rounded-lg px-2.5 py-1.5 focus:outline-none focus:border-indigo-500 font-medium"
          >
            {SCENARIOS.map(sc => (
              <option key={sc.id} value={sc.id}>
                {sc.title} (${sc.amount})
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Incident Context Banner */}
      <div className="p-4 rounded-xl bg-slate-900/60 border border-slate-800 flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-start gap-3">
          <div className="p-2 rounded-lg bg-amber-500/10 border border-amber-500/20 text-amber-400 shrink-0 mt-0.5">
            <AlertTriangle className="w-4 h-4" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <span className="text-sm font-semibold text-slate-200">Incident Stimulus:</span>
              <span className="font-mono text-xs text-amber-300">{selectedScenario.errorType}</span>
            </div>
            <p className="text-xs text-slate-400 mt-0.5">
              Reference: <span className="font-mono text-slate-300">{selectedScenario.merchantRef}</span> • Transaction Value: <span className="font-bold text-slate-200">${selectedScenario.amount} {selectedScenario.currency}</span>
            </p>
          </div>
        </div>
        <div className="flex items-center gap-2 text-xs text-slate-400 font-mono">
          <span className="px-2 py-1 rounded bg-slate-800 border border-slate-700/60 text-slate-300">
            PROVIDER: PAYPAL
          </span>
          <span className="px-2 py-1 rounded bg-slate-800 border border-slate-700/60 text-slate-300">
            ENV: SANDBOX
          </span>
        </div>
      </div>

      {/* Split Comparison Columns */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Left Column: With PayBridge AI */}
        <div className="rounded-2xl border border-emerald-500/30 bg-emerald-950/10 p-6 flex flex-col justify-between relative overflow-hidden shadow-lg shadow-emerald-950/20">
          <div className="absolute top-0 right-0 w-32 h-32 bg-emerald-500/5 rounded-full blur-2xl pointer-events-none" />
          
          <div>
            {/* Column Header */}
            <div className="flex items-center justify-between border-b border-emerald-800/40 pb-4 mb-5">
              <div className="flex items-center gap-2.5">
                <div className="p-2 rounded-lg bg-emerald-500/20 text-emerald-400 border border-emerald-500/30">
                  <ShieldCheck className="w-5 h-5" />
                </div>
                <div>
                  <h2 className="text-base font-bold text-emerald-300">With PayBridge AI</h2>
                  <p className="text-xs text-emerald-400/80">Deterministic Truth Resolution & Blind-Retry Guard</p>
                </div>
              </div>
              <span className="text-xs font-mono font-bold px-2.5 py-1 rounded-full bg-emerald-500/20 text-emerald-300 border border-emerald-500/30">
                ACTIVE
              </span>
            </div>

            {/* Step-by-Step Flow */}
            <div className="space-y-4">
              <div className="flex items-start gap-3 p-3 rounded-xl bg-slate-900/60 border border-emerald-900/40">
                <div className="w-6 h-6 rounded-full bg-emerald-500/20 text-emerald-400 flex items-center justify-center shrink-0 text-xs font-bold border border-emerald-500/30">
                  1
                </div>
                <div>
                  <h3 className="text-xs font-semibold text-slate-200">Payment Case Ingested (UNKNOWN State)</h3>
                  <p className="text-xs text-slate-400 mt-0.5">
                    Upon HTTP timeout, PayBridge creates a case with Provider State <code className="text-slate-300 bg-slate-800 px-1 py-0.5 rounded">UNKNOWN</code> instead of marking order as failed.
                  </p>
                </div>
              </div>

              <div className="flex items-start gap-3 p-3 rounded-xl bg-slate-900/60 border border-emerald-900/40">
                <div className="w-6 h-6 rounded-full bg-emerald-500/20 text-emerald-400 flex items-center justify-center shrink-0 text-xs font-bold border border-emerald-500/30">
                  2
                </div>
                <div>
                  <h3 className="text-xs font-semibold text-slate-200">Blind-Retry Defense Locked (409 Conflict)</h3>
                  <p className="text-xs text-slate-400 mt-0.5">
                    Acquires in-process and database atomic lock. If customer reloads or clicks "Pay Again", PayBridge returns <code>CAPTURE_IN_FLIGHT</code> to block double-billing.
                  </p>
                </div>
              </div>

              <div className="flex items-start gap-3 p-3 rounded-xl bg-slate-900/60 border border-emerald-900/40">
                <div className="w-6 h-6 rounded-full bg-emerald-500/20 text-emerald-400 flex items-center justify-center shrink-0 text-xs font-bold border border-emerald-500/30">
                  3
                </div>
                <div>
                  <h3 className="text-xs font-semibold text-slate-200">Authoritative Read-Only Probe</h3>
                  <p className="text-xs text-slate-400 mt-0.5">
                    Connector inspects provider capture endpoint directly without creating duplicate transactions. Confirms <span className="text-emerald-400 font-mono font-semibold">COMPLETED</span>.
                  </p>
                </div>
              </div>

              <div className="flex items-start gap-3 p-3 rounded-xl bg-slate-900/60 border border-emerald-900/40">
                <div className="w-6 h-6 rounded-full bg-emerald-500/20 text-emerald-400 flex items-center justify-center shrink-0 text-xs font-bold border border-emerald-500/30">
                  4
                </div>
                <div>
                  <h3 className="text-xs font-semibold text-slate-200">Safe Reader-Only Reconciliation</h3>
                  <p className="text-xs text-slate-400 mt-0.5">
                    Updates internal merchant record to <code className="text-emerald-400 bg-slate-800 px-1 py-0.5 rounded">SYNCHRONISED</code>, issues Payment Truth Passport with SHA-256 cryptographic seal.
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Bottom Financial Outcome */}
          <div className="mt-6 pt-4 border-t border-emerald-800/40 bg-emerald-950/30 -mx-6 -mb-6 p-6 rounded-b-2xl">
            <h3 className="text-xs font-bold text-emerald-400 uppercase tracking-wider mb-2">Net Business Result</h3>
            <div className="grid grid-cols-2 gap-3 text-xs">
              <div className="p-2.5 rounded-lg bg-slate-900/80 border border-emerald-900/50">
                <span className="text-slate-400 block text-[11px]">Customer Duplicate Charge</span>
                <span className="font-bold text-emerald-400 text-sm">$0.00 (Zero)</span>
              </div>
              <div className="p-2.5 rounded-lg bg-slate-900/80 border border-emerald-900/50">
                <span className="text-slate-400 block text-[11px]">Dispute / Chargeback Fees</span>
                <span className="font-bold text-emerald-400 text-sm">$0.00 Avoided</span>
              </div>
              <div className="p-2.5 rounded-lg bg-slate-900/80 border border-emerald-900/50">
                <span className="text-slate-400 block text-[11px]">Time to Full Reconciliation</span>
                <span className="font-bold text-slate-200 text-sm">~ 28 seconds</span>
              </div>
              <div className="p-2.5 rounded-lg bg-slate-900/80 border border-emerald-900/50">
                <span className="text-slate-400 block text-[11px]">Manual Engineering Review</span>
                <span className="font-bold text-slate-200 text-sm">0 minutes</span>
              </div>
            </div>
          </div>
        </div>

        {/* Right Column: Standard Default Behavior (Without PayBridge) */}
        <div className="rounded-2xl border border-rose-500/30 bg-rose-950/10 p-6 flex flex-col justify-between relative overflow-hidden shadow-lg shadow-rose-950/20">
          <div className="absolute top-0 right-0 w-32 h-32 bg-rose-500/5 rounded-full blur-2xl pointer-events-none" />

          <div>
            {/* Column Header */}
            <div className="flex items-center justify-between border-b border-rose-800/40 pb-4 mb-5">
              <div className="flex items-center gap-2.5">
                <div className="p-2 rounded-lg bg-rose-500/20 text-rose-400 border border-rose-500/30">
                  <ShieldAlert className="w-5 h-5" />
                </div>
                <div>
                  <h2 className="text-base font-bold text-rose-300">Without PayBridge (Industry Default)</h2>
                  <p className="text-xs text-rose-400/80">Blind Retry on Timeout & Uncoordinated States</p>
                </div>
              </div>
              <span className="text-xs font-mono font-bold px-2.5 py-1 rounded-full bg-rose-500/20 text-rose-300 border border-rose-500/30">
                DEFAULT
              </span>
            </div>

            {/* Step-by-Step Flow */}
            <div className="space-y-4">
              <div className="flex items-start gap-3 p-3 rounded-xl bg-slate-900/60 border border-rose-900/40">
                <div className="w-6 h-6 rounded-full bg-rose-500/20 text-rose-400 flex items-center justify-center shrink-0 text-xs font-bold border border-rose-500/30">
                  1
                </div>
                <div>
                  <h3 className="text-xs font-semibold text-slate-200">Timeout Treated as Outright Failure</h3>
                  <p className="text-xs text-slate-400 mt-0.5">
                    Merchant UI displays error banner: <span className="text-rose-400 italic">"Payment failed! Please try again or use another card."</span>
                  </p>
                </div>
              </div>

              <div className="flex items-start gap-3 p-3 rounded-xl bg-slate-900/60 border border-rose-900/40">
                <div className="w-6 h-6 rounded-full bg-rose-500/20 text-rose-400 flex items-center justify-center shrink-0 text-xs font-bold border border-rose-500/30">
                  2
                </div>
                <div>
                  <h3 className="text-xs font-semibold text-slate-200">Customer Retries Checkout (Blind Retry)</h3>
                  <p className="text-xs text-slate-400 mt-0.5">
                    Customer clicks "Pay Now" again. A second order/capture request is sent to the provider without cross-attempt deduplication.
                  </p>
                </div>
              </div>

              <div className="flex items-start gap-3 p-3 rounded-xl bg-slate-900/60 border border-rose-900/40">
                <div className="w-6 h-6 rounded-full bg-rose-500/20 text-rose-400 flex items-center justify-center shrink-0 text-xs font-bold border border-rose-500/30">
                  3
                </div>
                <div>
                  <h3 className="text-xs font-semibold text-slate-200">Customer Double-Billed</h3>
                  <p className="text-xs text-slate-400 mt-0.5">
                    Both transaction 1 and transaction 2 settle. Customer sees two identical debit charges of <span className="text-rose-400 font-bold">${selectedScenario.amount}</span> on their statement.
                  </p>
                </div>
              </div>

              <div className="flex items-start gap-3 p-3 rounded-xl bg-slate-900/60 border border-rose-900/40">
                <div className="w-6 h-6 rounded-full bg-rose-500/20 text-rose-400 flex items-center justify-center shrink-0 text-xs font-bold border border-rose-500/30">
                  4
                </div>
                <div>
                  <h3 className="text-xs font-semibold text-slate-200">Chargeback Dispute & Churn</h3>
                  <p className="text-xs text-slate-400 mt-0.5">
                    Customer files dispute with bank. Provider levies non-refundable <span className="text-rose-400 font-semibold">$15.00 dispute penalty</span>. Merchant trust is severely compromised.
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Bottom Financial Outcome */}
          <div className="mt-6 pt-4 border-t border-rose-800/40 bg-rose-950/30 -mx-6 -mb-6 p-6 rounded-b-2xl">
            <h3 className="text-xs font-bold text-rose-400 uppercase tracking-wider mb-2">Net Business Damage</h3>
            <div className="grid grid-cols-2 gap-3 text-xs">
              <div className="p-2.5 rounded-lg bg-slate-900/80 border border-rose-900/50">
                <span className="text-slate-400 block text-[11px]">Customer Duplicate Charge</span>
                <span className="font-bold text-rose-400 text-sm">+${selectedScenario.amount} USD</span>
              </div>
              <div className="p-2.5 rounded-lg bg-slate-900/80 border border-rose-900/50">
                <span className="text-slate-400 block text-[11px]">Processor Dispute Penalty</span>
                <span className="font-bold text-rose-400 text-sm">-$15.00 per claim</span>
              </div>
              <div className="p-2.5 rounded-lg bg-slate-900/80 border border-rose-900/50">
                <span className="text-slate-400 block text-[11px]">Dispute Settlement Window</span>
                <span className="font-bold text-slate-200 text-sm">14 – 30 days</span>
              </div>
              <div className="p-2.5 rounded-lg bg-slate-900/80 border border-rose-900/50">
                <span className="text-slate-400 block text-[11px]">Support / Ops Investigation</span>
                <span className="font-bold text-slate-200 text-sm">~ 45 minutes / ticket</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Comparative ROI & Performance Matrix */}
      <div className="rounded-2xl border border-slate-800 bg-slate-900/50 p-6">
        <h3 className="text-sm font-bold text-slate-200 uppercase tracking-wider mb-4 flex items-center gap-2">
          <Scale className="w-4 h-4 text-indigo-400" />
          Direct Technical Comparison Matrix
        </h3>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead>
              <tr className="border-b border-slate-800 text-slate-400 uppercase font-mono">
                <th className="py-2.5 px-3">Evaluation Dimension</th>
                <th className="py-2.5 px-3 text-emerald-400 font-bold">PayBridge AI Architecture</th>
                <th className="py-2.5 px-3 text-rose-400 font-bold">Industry Standard Fallback</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/60 font-medium">
              <tr>
                <td className="py-3 px-3 text-slate-300 font-semibold">State Inconsistency Representation</td>
                <td className="py-3 px-3 text-emerald-300">
                  Explicit 3D State Vector (<code className="bg-slate-800 px-1 py-0.5 rounded text-emerald-400">UNKNOWN</code> provider state, never faked)
                </td>
                <td className="py-3 px-3 text-rose-300">
                  Binary boolean <code className="bg-slate-800 px-1 py-0.5 rounded text-rose-400">is_paid: false</code>
                </td>
              </tr>
              <tr>
                <td className="py-3 px-3 text-slate-300 font-semibold">Concurrency & Duplicate Lock</td>
                <td className="py-3 px-3 text-emerald-300">
                  In-process mutual exclusion + DB conditional status lock (returns 409 Conflict)
                </td>
                <td className="py-3 px-3 text-rose-300">
                  No capture mutex; every button click dispatches raw provider request
                </td>
              </tr>
              <tr>
                <td className="py-3 px-3 text-slate-300 font-semibold">Reconciliation Behavior</td>
                <td className="py-3 px-3 text-emerald-300">
                  Reader-only authoritative polling (zero payment mutations)
                </td>
                <td className="py-3 px-3 text-rose-300">
                  Blind retry or destructive force-cancellation
                </td>
              </tr>
              <tr>
                <td className="py-3 px-3 text-slate-300 font-semibold">Audit & Dispute Defensibility</td>
                <td className="py-3 px-3 text-emerald-300">
                  Immutable Tier 1–6 evidence store with SHA-256 payload seal
                </td>
                <td className="py-3 px-3 text-rose-300">
                  Uncorrelated ephemeral application log files
                </td>
              </tr>
              <tr>
                <td className="py-3 px-3 text-slate-300 font-semibold">Agent Role Safety Boundary</td>
                <td className="py-3 px-3 text-emerald-300">
                  AI strictly limited to read/classify/recommend. Cannot mutate money.
                </td>
                <td className="py-3 px-3 text-rose-300">
                  Unconstrained autonomous tool calls or ad-hoc admin actions
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
