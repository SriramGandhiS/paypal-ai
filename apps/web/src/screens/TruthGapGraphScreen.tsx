import React, { useState } from 'react';
import { 
  Network, 
  AlertTriangle, 
  CheckCircle2, 
  ShieldCheck, 
  Clock, 
  ArrowRight,
  Info
} from 'lucide-react';
import { StatusPill } from '../components/StatusPill';

interface GraphNode {
  id: string;
  title: string;
  subtitle: string;
  status: 'SYNCHRONISED' | 'TIMEOUT' | 'CONFIRMED' | 'BRIDGED';
  details: string;
  telemetry: Record<string, string>;
}

export const TruthGapGraphScreen: React.FC = () => {
  const [selectedNode, setSelectedNode] = useState<string>('gap');

  const nodes: Record<string, GraphNode> = {
    sender: {
      id: 'sender',
      title: 'Customer Intent',
      subtitle: 'Buyer Checkout Window',
      status: 'CONFIRMED',
      details: 'Customer initiated PayPal checkout and approved order ORDER-SANDBOX-1794316132 for $149.50 USD.',
      telemetry: {
        'Buyer Action': 'APPROVED',
        'Amount': '$149.50 USD',
        'Time': '18:38:52 UTC',
      },
    },
    merchant: {
      id: 'merchant',
      title: 'Merchant Website',
      subtitle: 'Internal Order Store',
      status: 'TIMEOUT',
      details: 'Merchant application server timed out while awaiting webhook. Order INV-2026-8941 remained in UNKNOWN/PENDING state.',
      telemetry: {
        'HTTP Response': '504 Gateway Timeout',
        'Internal State': 'INCONSISTENT',
        'Blind-Retry Danger': 'CRITICAL (User would be double-charged)',
      },
    },
    gap: {
      id: 'gap',
      title: 'The Truth Gap',
      subtitle: 'Asymmetric Information severed',
      status: 'BRIDGED',
      details: 'The critical point of failure in modern e-commerce: PayPal completed capture, but merchant order system has no proof.',
      telemetry: {
        'Disconnect Duration': '15.2 seconds',
        'Gap Resolution': 'PayBridge AI Reader Query',
        'Financial Truth': 'Secured',
      },
    },
    paypal: {
      id: 'paypal',
      title: 'PayPal Sandbox',
      subtitle: 'Payment Rail Processor',
      status: 'CONFIRMED',
      details: 'Authoritative provider capture CAP-PAYPAL-782199 completed successfully. Funds debited from buyer and credited to seller.',
      telemetry: {
        'Capture ID': 'CAP-PAYPAL-782199',
        'Provider State': 'CAPTURE_COMPLETED',
        'Seller Protection': 'ELIGIBLE',
      },
    },
    paybridge: {
      id: 'paybridge',
      title: 'PayBridge Engine',
      subtitle: 'Deterministic Truth Vector',
      status: 'SYNCHRONISED',
      details: 'Ingested Tier 1 capture proof, blocked duplicate blind retries, incremented state vector v3, and issued Truth Passport.',
      telemetry: {
        'Vector State': 'SYNCHRONISED',
        'Evidence Hash': '750617fa580caf88...',
        'Blind Retry': 'BLOCKED',
      },
    },
  };

  const active = nodes[selectedNode] || nodes.gap;

  return (
    <div className="space-y-8 animate-fade-in max-w-5xl mx-auto pb-12">
      <div>
        <div className="flex items-center gap-2">
          <Network className="w-6 h-6 text-cyan-400" />
          <h1 className="text-2xl font-bold tracking-tight text-slate-100">
            Truth Gap Graph
          </h1>
        </div>
        <p className="text-sm text-slate-400 mt-1">
          Visual topology of information asymmetry and automated bridge reconciliation.
        </p>
      </div>

      {/* Interactive Topology Visualizer */}
      <div className="p-6 sm:p-8 rounded-2xl bg-slate-900/80 border border-slate-800 shadow-xl space-y-6">
        <div className="flex items-center justify-between">
          <span className="text-xs font-mono uppercase text-slate-400 tracking-wider font-bold">
            Interactive Reconciliation Topology (Click any node to inspect)
          </span>
          <span className="text-xs text-indigo-400 font-mono">
            Active Selection: {active.title}
          </span>
        </div>

        {/* Node Flow Diagram */}
        <div className="grid grid-cols-1 md:grid-cols-5 gap-3 items-center py-4">
          {/* Node 1: Customer */}
          <div 
            onClick={() => setSelectedNode('sender')}
            className={`p-4 rounded-xl border text-center transition-all cursor-pointer ${
              selectedNode === 'sender'
                ? 'bg-indigo-950/70 border-indigo-500 shadow-md shadow-indigo-500/20'
                : 'bg-slate-950/60 border-slate-800 hover:border-slate-700'
            }`}
          >
            <div className="w-2.5 h-2.5 rounded-full bg-emerald-400 mx-auto mb-2"></div>
            <div className="text-xs font-bold text-slate-100">1. Customer</div>
            <div className="text-[11px] text-slate-400 mt-0.5">Approved Order</div>
          </div>

          {/* Node 2: Merchant (Disconnect) */}
          <div 
            onClick={() => setSelectedNode('merchant')}
            className={`p-4 rounded-xl border text-center transition-all cursor-pointer ${
              selectedNode === 'merchant'
                ? 'bg-rose-950/70 border-rose-500 shadow-md shadow-rose-500/20'
                : 'bg-slate-950/60 border-slate-800 hover:border-slate-700'
            }`}
          >
            <div className="w-2.5 h-2.5 rounded-full bg-rose-400 mx-auto mb-2"></div>
            <div className="text-xs font-bold text-slate-100">2. Merchant Store</div>
            <div className="text-[11px] text-rose-300 mt-0.5">Timeout Disconnect</div>
          </div>

          {/* Node 3: The Gap */}
          <div 
            onClick={() => setSelectedNode('gap')}
            className={`p-4 rounded-xl border text-center transition-all cursor-pointer relative overflow-hidden ${
              selectedNode === 'gap'
                ? 'bg-amber-950/70 border-amber-500 shadow-md shadow-amber-500/20 ring-1 ring-amber-500'
                : 'bg-amber-950/30 border-amber-800/60 hover:border-amber-600'
            }`}
          >
            <div className="w-2.5 h-2.5 rounded-full bg-amber-400 mx-auto mb-2 animate-ping"></div>
            <div className="text-xs font-bold text-amber-200">3. TRUTH GAP</div>
            <div className="text-[11px] text-amber-300/80 mt-0.5">Asymmetry Detected</div>
          </div>

          {/* Node 4: PayPal Provider */}
          <div 
            onClick={() => setSelectedNode('paypal')}
            className={`p-4 rounded-xl border text-center transition-all cursor-pointer ${
              selectedNode === 'paypal'
                ? 'bg-indigo-950/70 border-indigo-500 shadow-md shadow-indigo-500/20'
                : 'bg-slate-950/60 border-slate-800 hover:border-slate-700'
            }`}
          >
            <div className="w-2.5 h-2.5 rounded-full bg-emerald-400 mx-auto mb-2"></div>
            <div className="text-xs font-bold text-slate-100">4. PayPal Rail</div>
            <div className="text-[11px] text-emerald-300 mt-0.5">Capture Success</div>
          </div>

          {/* Node 5: PayBridge Reconciled */}
          <div 
            onClick={() => setSelectedNode('paybridge')}
            className={`p-4 rounded-xl border text-center transition-all cursor-pointer ${
              selectedNode === 'paybridge'
                ? 'bg-cyan-950/70 border-cyan-500 shadow-md shadow-cyan-500/20 ring-1 ring-cyan-500'
                : 'bg-slate-950/60 border-slate-800 hover:border-slate-700'
            }`}
          >
            <div className="w-2.5 h-2.5 rounded-full bg-cyan-400 mx-auto mb-2"></div>
            <div className="text-xs font-bold text-slate-100">5. PayBridge AI</div>
            <div className="text-[11px] text-cyan-300 mt-0.5">Truth Reconciled</div>
          </div>
        </div>

        {/* Selected Node Inspector Card */}
        <div className="p-5 rounded-xl bg-slate-950/90 border border-slate-800 space-y-4">
          <div className="flex items-center justify-between border-b border-slate-800 pb-3">
            <div>
              <span className="text-xs font-mono text-slate-400 uppercase font-bold">Node Telemetry Inspector</span>
              <h3 className="text-base font-bold text-slate-100 mt-0.5">{active.title} — {active.subtitle}</h3>
            </div>
            <StatusPill status={active.status} />
          </div>

          <p className="text-sm text-slate-300 leading-relaxed">
            {active.details}
          </p>

          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 pt-2">
            {Object.entries(active.telemetry).map(([k, v]) => (
              <div key={k} className="p-3 rounded-lg bg-slate-900 border border-slate-800">
                <span className="text-[10px] font-mono text-slate-400 uppercase block">{k}</span>
                <span className="text-xs font-mono font-semibold text-slate-200 mt-1 block truncate">
                  {v}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};
