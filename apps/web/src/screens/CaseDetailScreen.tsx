import React, { useState, useEffect } from 'react';
import { 
  ArrowLeft, 
  ShieldCheck, 
  Clock, 
  FileCheck, 
  AlertOctagon, 
  CheckCircle2, 
  ExternalLink,
  Lock,
  ChevronDown,
  ChevronUp,
  Award,
  Network,
  RotateCcw
} from 'lucide-react';
import { api } from '../api/client';
import { Case, Attempt, EvidenceItem, SafeAction } from '../api/types';
import { StatusPill } from '../components/StatusPill';

interface CaseDetailScreenProps {
  caseId: string;
  onBack: () => void;
  onNavigate: (tab: string) => void;
}

export const CaseDetailScreen: React.FC<CaseDetailScreenProps> = ({
  caseId,
  onBack,
  onNavigate,
}) => {
  const [caseData, setCaseData] = useState<Case | null>(null);
  const [attempts, setAttempts] = useState<Attempt[]>([]);
  const [evidence, setEvidence] = useState<EvidenceItem[]>([]);
  const [actions, setActions] = useState<SafeAction[]>([]);
  const [loading, setLoading] = useState(true);
  const [expandedPayload, setExpandedPayload] = useState<string | null>(null);
  const [actionExecuted, setActionExecuted] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      const c = await api.getCase(caseId);
      setCaseData(c);
      if (c) {
        const atts = await api.listAttempts(c.id);
        const evs = await api.listEvidence(c.id);
        const acts = await api.getSafeActions(c.id);
        setAttempts(atts);
        setEvidence(evs);
        setActions(acts);
      }
      setLoading(false);
    };
    fetchData();
  }, [caseId]);

  const handleExecuteSafeAction = (action: SafeAction) => {
    if (!caseData) return;
    // Execute safe reconciliation
    api.updateCaseState(caseData.id, {
      application_state: 'SYNCHRONISED',
      case_status: 'RESOLVED',
      state_version: caseData.state_version + 1,
    });
    setCaseData({
      ...caseData,
      application_state: 'SYNCHRONISED',
      case_status: 'RESOLVED',
      state_version: caseData.state_version + 1,
    });
    setActionExecuted(true);
  };

  if (loading || !caseData) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="w-8 h-8 rounded-full border-2 border-indigo-500 border-t-transparent animate-spin"></div>
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-fade-in pb-12">
      {/* Top Breadcrumb & Actions */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <button
          onClick={onBack}
          className="flex items-center gap-2 text-sm text-slate-400 hover:text-slate-200 cursor-pointer font-medium"
        >
          <ArrowLeft className="w-4 h-4" />
          Back to Cases
        </button>
        <div className="flex items-center gap-3">
          <button
            onClick={() => onNavigate('passport')}
            className="flex items-center gap-2 px-3.5 py-1.5 bg-indigo-600/20 hover:bg-indigo-600/30 text-indigo-300 border border-indigo-500/30 rounded-lg text-xs font-medium cursor-pointer transition-all"
          >
            <Award className="w-4 h-4 text-indigo-400" />
            View Truth Passport
          </button>
          <button
            onClick={() => onNavigate('graph')}
            className="flex items-center gap-2 px-3.5 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 border border-slate-700 rounded-lg text-xs font-medium cursor-pointer transition-all"
          >
            <Network className="w-4 h-4 text-cyan-400" />
            Inspect Gap Graph
          </button>
        </div>
      </div>

      {/* Case Header Banner */}
      <div className="p-6 rounded-2xl bg-slate-900/80 border border-slate-800 shadow-sm">
        <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-6">
          <div>
            <div className="flex items-center gap-3">
              <span className="font-mono text-2xl font-bold text-slate-100 tracking-tight">
                {caseData.external_reference}
              </span>
              <span className="text-xs px-2.5 py-0.5 rounded font-mono bg-slate-800 text-slate-400 border border-slate-700">
                UUID: {caseData.id.slice(0, 18)}...
              </span>
            </div>
            <div className="flex flex-wrap items-center gap-4 text-xs text-slate-400 mt-2">
              <span>Purpose: <strong className="text-slate-200">{caseData.payment_purpose}</strong></span>
              <span>•</span>
              <span>Corridor: <strong className="text-slate-200">{caseData.sender_country} → {caseData.receiver_country}</strong></span>
              <span>•</span>
              <span>Environment: <strong className="text-slate-200">{caseData.environment}</strong></span>
              <span>•</span>
              <span>Created: {new Date(caseData.created_at).toLocaleString()}</span>
            </div>
          </div>

          <div className="flex items-baseline gap-2 text-right">
            <span className="text-xs text-slate-400 font-medium">Expected:</span>
            <span className="text-3xl font-extrabold text-slate-100 font-mono tracking-tight tabular-nums">
              ${caseData.amount}
            </span>
            <span className="text-sm font-semibold text-slate-400">{caseData.currency}</span>
          </div>
        </div>

        {/* 3D State Vector Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6 pt-6 border-t border-slate-800/80">
          <div className="p-4 rounded-xl bg-slate-950/60 border border-slate-800/80">
            <span className="text-xs font-mono uppercase text-slate-400">1. Provider Truth State</span>
            <div className="mt-2 flex items-center justify-between">
              <StatusPill status={caseData.provider_state} type="provider" />
              <span className="text-[10px] font-mono text-slate-400">Authoritative GET</span>
            </div>
            <p className="text-xs text-slate-400 mt-2">
              Ground truth reported directly by PayPal Sandbox capture API.
            </p>
          </div>

          <div className="p-4 rounded-xl bg-slate-950/60 border border-slate-800/80">
            <span className="text-xs font-mono uppercase text-slate-400">2. Application State</span>
            <div className="mt-2 flex items-center justify-between">
              <StatusPill status={caseData.application_state} type="app" />
              <span className="text-[10px] font-mono text-slate-400">Merchant Sync</span>
            </div>
            <p className="text-xs text-slate-400 mt-2">
              Internal invoice status within the merchant order management system.
            </p>
          </div>

          <div className="p-4 rounded-xl bg-slate-950/60 border border-slate-800/80">
            <span className="text-xs font-mono uppercase text-slate-400">3. Master Case Status</span>
            <div className="mt-2 flex items-center justify-between">
              <StatusPill status={caseData.case_status} type="case" />
              <span className="text-[10px] font-mono text-slate-400">Vector v{caseData.state_version}</span>
            </div>
            <p className="text-xs text-slate-400 mt-2">
              Reconciliation status determining whether human review or safe action is required.
            </p>
          </div>
        </div>
      </div>

      {/* Next Safe Action Recommendation Engine */}
      <div className="p-6 rounded-2xl bg-gradient-to-r from-indigo-950/30 via-slate-900 to-slate-900 border border-indigo-500/30 shadow-lg shadow-indigo-500/5">
        <div className="flex items-center gap-2 text-indigo-400 text-xs font-mono uppercase font-bold tracking-wider">
          <ShieldCheck className="w-4 h-4 text-indigo-400" />
          Next Safe Action Engine (Deterministic Policy Guard)
        </div>

        {actionExecuted ? (
          <div className="mt-4 p-4 rounded-xl bg-emerald-950/40 border border-emerald-800/60 flex items-center gap-3">
            <CheckCircle2 className="w-5 h-5 text-emerald-400 flex-shrink-0" />
            <div>
              <h4 className="text-sm font-semibold text-emerald-200">Safe Action Executed Successfully</h4>
              <p className="text-xs text-slate-300 mt-0.5">
                The order record has been reconciled with authoritative provider ground truth. Blind retries are disarmed.
              </p>
            </div>
          </div>
        ) : actions.length > 0 ? (
          <div className="mt-4 space-y-3">
            {actions.map((act) => (
              <div 
                key={act.id}
                className="p-4 rounded-xl bg-slate-950/60 border border-slate-800 flex flex-col sm:flex-row sm:items-center justify-between gap-4"
              >
                <div>
                  <div className="flex items-center gap-2">
                    <span className="font-semibold text-sm text-slate-100">{act.title}</span>
                    <StatusPill status={act.risk_level} />
                  </div>
                  <p className="text-xs text-slate-300 mt-1">{act.description}</p>
                  <p className="text-xs text-amber-300/90 mt-1 font-mono">
                    ⚠️ Counterfactual: {act.counterfactual_warning}
                  </p>
                </div>
                <button
                  onClick={() => handleExecuteSafeAction(act)}
                  className="flex-shrink-0 px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white font-medium text-xs rounded-lg shadow-md transition-all cursor-pointer"
                >
                  Execute Safe Action
                </button>
              </div>
            ))}
          </div>
        ) : null}
      </div>

      {/* Two Column Layout: Attempts + Evidence */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Left: Payment Attempts Timeline */}
        <div className="rounded-xl bg-slate-900/60 border border-slate-800 p-6 space-y-4">
          <div className="flex items-center justify-between border-b border-slate-800 pb-3">
            <h3 className="font-semibold text-sm text-slate-100 flex items-center gap-2">
              <Clock className="w-4 h-4 text-indigo-400" />
              Payment Attempts (Blind-Retry Guard)
            </h3>
            <span className="text-xs font-mono text-slate-400">
              {attempts.length} Attempt(s)
            </span>
          </div>

          {attempts.length === 0 ? (
            <p className="text-xs text-slate-500 py-6 text-center">No payment attempts initiated yet.</p>
          ) : (
            <div className="space-y-3">
              {attempts.map((att) => (
                <div key={att.id} className="p-3.5 rounded-lg bg-slate-950/50 border border-slate-800 space-y-2">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <span className="text-xs font-bold font-mono text-indigo-400">Attempt #{att.attempt_no}</span>
                      <StatusPill status={att.status} />
                    </div>
                    <span className="text-[11px] font-mono text-slate-500">
                      {new Date(att.created_at).toLocaleTimeString()}
                    </span>
                  </div>
                  <div className="text-xs font-mono text-slate-400 space-y-0.5">
                    <div>Order Op ID: <span className="text-slate-200">{att.order_operation_request_id}</span></div>
                    {att.capture_operation_request_id && (
                      <div>Capture Op ID: <span className="text-slate-200">{att.capture_operation_request_id}</span></div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Right: Evidence Provenance Store */}
        <div className="rounded-xl bg-slate-900/60 border border-slate-800 p-6 space-y-4">
          <div className="flex items-center justify-between border-b border-slate-800 pb-3">
            <h3 className="font-semibold text-sm text-slate-100 flex items-center gap-2">
              <FileCheck className="w-4 h-4 text-cyan-400" />
              Immutable Evidence Trail (SHA-256 Verified)
            </h3>
            <span className="text-xs font-mono text-slate-400">
              {evidence.length} Evidence Records
            </span>
          </div>

          {evidence.length === 0 ? (
            <p className="text-xs text-slate-500 py-6 text-center">No evidence ingested for this case.</p>
          ) : (
            <div className="space-y-3">
              {evidence.map((evi) => {
                const isExpanded = expandedPayload === evi.id;
                return (
                  <div key={evi.id} className="p-3.5 rounded-lg bg-slate-950/50 border border-slate-800 space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-xs font-bold font-mono text-slate-200">{evi.kind.replace(/_/g, ' ')}</span>
                      <StatusPill status={`TIER ${evi.tier} — ${evi.source_trust}`} />
                    </div>

                    <div className="flex items-center justify-between text-[11px] font-mono text-slate-400">
                      <span>SHA-256: {evi.payload_sha256.slice(0, 16)}...</span>
                      <button
                        onClick={() => setExpandedPayload(isExpanded ? null : evi.id)}
                        className="text-indigo-400 hover:text-indigo-300 flex items-center gap-1 cursor-pointer"
                      >
                        {isExpanded ? 'Hide Payload' : 'Inspect Payload'}
                        {isExpanded ? <ChevronUp className="w-3 h-3" /> : <ChevronDown className="w-3 h-3" />}
                      </button>
                    </div>

                    {isExpanded && (
                      <pre className="p-3 rounded bg-slate-900 border border-slate-800 text-[11px] font-mono text-slate-300 overflow-x-auto mt-2">
                        {JSON.stringify(evi.redacted_extract, null, 2)}
                      </pre>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
