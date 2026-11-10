import React, { useState } from 'react';
import { 
  PlayCircle, 
  RotateCcw, 
  CheckCircle2, 
  Clock, 
  ArrowRight, 
  ShieldCheck, 
  AlertTriangle, 
  FileText, 
  Lock, 
  Server, 
  Search, 
  Award,
  ChevronRight,
  Code
} from 'lucide-react';
import { api } from '../api/client';
import { StatusPill } from '../components/StatusPill';

interface DemoStep {
  step: number;
  title: string;
  category: 'SETUP' | 'CHECKOUT' | 'APPROVAL' | 'INCIDENT' | 'PROBE' | 'RECONCILE' | 'SEAL';
  description: string;
  endpoint: string;
  httpMethod: string;
  requestPayload?: any;
  responsePayload?: any;
  stateVectorAfter: {
    provider_state: string;
    application_state: string;
    case_status: string;
  };
}

const DEMO_STEPS: DemoStep[] = [
  {
    step: 1,
    title: 'Initialize Payment Case',
    category: 'SETUP',
    description: 'Merchant initializes a case container to track customer intent and financial boundaries.',
    endpoint: '/api/v1/cases',
    httpMethod: 'POST',
    requestPayload: {
      external_reference: 'INV-2026-8941',
      amount_cents: 14950,
      currency: 'USD',
      sender_country: 'US',
      receiver_country: 'DE',
      payment_purpose: 'SAAS_SUBSCRIPTION',
    },
    responsePayload: {
      id: 'c101-demo-4567-89ab-cdef00000001',
      external_reference: 'INV-2026-8941',
      amount: '149.50',
      currency: 'USD',
      provider_state: 'UNKNOWN',
      application_state: 'INCONSISTENT',
      case_status: 'OBSERVING',
      state_version: 1,
    },
    stateVectorAfter: {
      provider_state: 'UNKNOWN',
      application_state: 'INCONSISTENT',
      case_status: 'OBSERVING',
    },
  },
  {
    step: 2,
    title: 'Initiate Checkout Attempt',
    category: 'CHECKOUT',
    description: 'PayBridge registers Attempt #1 with unique operation request ID to guard idempotency.',
    endpoint: '/api/v1/cases/c101/checkout/paypal/create',
    httpMethod: 'POST',
    requestPayload: {
      case_id: 'c101-demo-4567-89ab-cdef00000001',
      return_url: 'https://merchant.example.com/checkout/return',
      cancel_url: 'https://merchant.example.com/checkout/cancel',
    },
    responsePayload: {
      attempt_id: 'att-101-1',
      provider_order_id: 'ORDER-SANDBOX-1794316132',
      status: 'ORDER_CREATED',
      approval_link: 'https://www.sandbox.paypal.com/checkoutnow?token=ORDER-SANDBOX-1794316132',
    },
    stateVectorAfter: {
      provider_state: 'UNKNOWN',
      application_state: 'INCONSISTENT',
      case_status: 'OBSERVING',
    },
  },
  {
    step: 3,
    title: 'Customer Approval in Sandbox',
    category: 'APPROVAL',
    description: 'Buyer authorizes payment inside PayPal sandbox checkout frame. Funds are ready for capture.',
    endpoint: '/api/v1/cases/c101/checkout/paypal/approve',
    httpMethod: 'POST',
    requestPayload: {
      provider_order_id: 'ORDER-SANDBOX-1794316132',
      payer_id: 'PAYER-SANDBOX-US-001',
    },
    responsePayload: {
      attempt_id: 'att-101-1',
      status: 'ORDER_APPROVED',
      captured: false,
    },
    stateVectorAfter: {
      provider_state: 'UNKNOWN',
      application_state: 'INCONSISTENT',
      case_status: 'OBSERVING',
    },
  },
  {
    step: 4,
    title: 'Gateway Timeout & Blind-Retry Lock',
    category: 'INCIDENT',
    description: 'Upstream gateway drops connection during capture execution. Mutex lock acquires to defend against duplicate charges.',
    endpoint: '/api/v1/cases/c101/checkout/paypal/capture',
    httpMethod: 'POST',
    requestPayload: {
      provider_order_id: 'ORDER-SANDBOX-1794316132',
      idempotency_key: 'idem-capture-001',
    },
    responsePayload: {
      error: 'UPSTREAM_GATEWAY_TIMEOUT',
      detail: 'Connection dropped after capture dispatch. Acquired in-process capture lock to prevent blind retry.',
      status: 'CAPTURE_IN_FLIGHT',
    },
    stateVectorAfter: {
      provider_state: 'UNKNOWN',
      application_state: 'API_TIMEOUT',
      case_status: 'NEEDS_HUMAN_REVIEW',
    },
  },
  {
    step: 5,
    title: 'Authoritative Provider Capture Query',
    category: 'PROBE',
    description: 'PayBridge Connector executes read-only inquiry directly against PayPal GET /v2/payments/captures.',
    endpoint: '/v2/payments/captures/CAP-PAYPAL-782199',
    httpMethod: 'GET',
    requestPayload: null,
    responsePayload: {
      id: 'CAP-PAYPAL-782199',
      status: 'COMPLETED',
      amount: { value: '149.50', currency_code: 'USD' },
      seller_protection: { status: 'ELIGIBLE' },
      create_time: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
    },
    stateVectorAfter: {
      provider_state: 'CAPTURE_COMPLETED',
      application_state: 'API_TIMEOUT',
      case_status: 'NEEDS_HUMAN_REVIEW',
    },
  },
  {
    step: 6,
    title: 'Reader-Only Truth Reconciliation',
    category: 'RECONCILE',
    description: 'Reconciler processes Tier 1 evidence, reconciles merchant order status without mutating payment rail, releases mutex.',
    endpoint: '/api/v1/cases/c101/reconcile',
    httpMethod: 'POST',
    requestPayload: {
      action: 'RECONCILE_MERCHANT_RECORD',
      confirmed_capture_id: 'CAP-PAYPAL-782199',
    },
    responsePayload: {
      reconciled: true,
      provider_state: 'CAPTURE_COMPLETED',
      application_state: 'SYNCHRONISED',
      case_status: 'RESOLVED',
      duplicate_charge_prevented: true,
    },
    stateVectorAfter: {
      provider_state: 'CAPTURE_COMPLETED',
      application_state: 'SYNCHRONISED',
      case_status: 'RESOLVED',
    },
  },
  {
    step: 7,
    title: 'Cryptographic Truth Passport Sealed',
    category: 'SEAL',
    description: 'Generates immutable Payment Truth Passport with SHA-256 seal and complete audit lineage for dispute immunity.',
    endpoint: '/api/v1/cases/c101/passport',
    httpMethod: 'GET',
    requestPayload: null,
    responsePayload: {
      passport_id: 'PASS-2026-c101-01',
      passport_hash: '4fc20f0e89261a829b68f023fcc582934dae262520ded004aa12ec660191ac6b',
      verification_status: 'VERIFIED_AUTHORITATIVE',
      evidence_tiers_included: [1, 4],
      duplicate_attempts_blocked: 1,
    },
    stateVectorAfter: {
      provider_state: 'CAPTURE_COMPLETED',
      application_state: 'SYNCHRONISED',
      case_status: 'RESOLVED',
    },
  },
];

interface GoldenPathDemoScreenProps {
  onNavigate: (tab: string) => void;
  onSelectCase: (caseId: string) => void;
}

export const GoldenPathDemoScreen: React.FC<GoldenPathDemoScreenProps> = ({
  onNavigate,
  onSelectCase,
}) => {
  const [currentStepIndex, setCurrentStepIndex] = useState<number>(0);
  const [isRunningAll, setIsRunningAll] = useState<boolean>(false);
  const [activeTab, setActiveTab] = useState<'response' | 'request'>('response');

  const currentStep = DEMO_STEPS[currentStepIndex];

  const handleNextStep = () => {
    if (currentStepIndex < DEMO_STEPS.length - 1) {
      const nextIdx = currentStepIndex + 1;
      setCurrentStepIndex(nextIdx);
      // Synchronize with API client state
      api.updateCaseState('c101-demo-4567-89ab-cdef00000001', DEMO_STEPS[nextIdx].stateVectorAfter);
    }
  };

  const handleReset = () => {
    setCurrentStepIndex(0);
    api.updateCaseState('c101-demo-4567-89ab-cdef00000001', DEMO_STEPS[0].stateVectorAfter);
  };

  const handleRunAll = async () => {
    setIsRunningAll(true);
    for (let i = currentStepIndex; i < DEMO_STEPS.length; i++) {
      setCurrentStepIndex(i);
      api.updateCaseState('c101-demo-4567-89ab-cdef00000001', DEMO_STEPS[i].stateVectorAfter);
      await new Promise(r => setTimeout(r, 900));
    }
    setIsRunningAll(false);
  };

  return (
    <div className="space-y-8 animate-fade-in max-w-6xl mx-auto pb-12">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <PlayCircle className="w-6 h-6 text-indigo-400" />
            <h1 className="text-2xl font-bold tracking-tight text-slate-100">
              Interactive Golden Path Demo
            </h1>
          </div>
          <p className="text-sm text-slate-400 mt-1">
            Step-by-step verification walkthrough: from gateway timeout incident to cryptographic Truth Passport sealing.
          </p>
        </div>

        {/* Action Controls */}
        <div className="flex items-center gap-2">
          <button
            onClick={handleReset}
            disabled={isRunningAll || currentStepIndex === 0}
            className="flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-semibold bg-slate-900 border border-slate-800 text-slate-300 hover:bg-slate-800 disabled:opacity-40 disabled:cursor-not-allowed transition-all"
          >
            <RotateCcw className="w-3.5 h-3.5" />
            Reset
          </button>
          <button
            onClick={handleNextStep}
            disabled={isRunningAll || currentStepIndex === DEMO_STEPS.length - 1}
            className="flex items-center gap-1.5 px-4 py-2 rounded-lg text-xs font-semibold bg-indigo-600 hover:bg-indigo-500 text-white disabled:opacity-40 disabled:cursor-not-allowed shadow-md shadow-indigo-600/20 transition-all"
          >
            <span>Next Step</span>
            <ChevronRight className="w-3.5 h-3.5" />
          </button>
          <button
            onClick={handleRunAll}
            disabled={isRunningAll || currentStepIndex === DEMO_STEPS.length - 1}
            className="flex items-center gap-1.5 px-4 py-2 rounded-lg text-xs font-semibold bg-cyan-600 hover:bg-cyan-500 text-white disabled:opacity-40 disabled:cursor-not-allowed shadow-md shadow-cyan-600/20 transition-all"
          >
            <PlayCircle className="w-3.5 h-3.5" />
            <span>{isRunningAll ? 'Running Simulation...' : 'Run All Steps'}</span>
          </button>
        </div>
      </div>

      {/* Stepper Navigation Bar */}
      <div className="rounded-2xl border border-slate-800 bg-slate-900/60 p-4">
        <div className="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-7 gap-2">
          {DEMO_STEPS.map((step, idx) => {
            const isCompleted = idx < currentStepIndex;
            const isCurrent = idx === currentStepIndex;

            return (
              <button
                key={step.step}
                onClick={() => {
                  if (!isRunningAll) {
                    setCurrentStepIndex(idx);
                    api.updateCaseState('c101-demo-4567-89ab-cdef00000001', step.stateVectorAfter);
                  }
                }}
                className={`flex flex-col items-start p-2.5 rounded-xl border text-left transition-all cursor-pointer ${
                  isCurrent
                    ? 'bg-indigo-950/60 border-indigo-500 text-indigo-300 shadow-sm'
                    : isCompleted
                    ? 'bg-slate-900/90 border-emerald-900/60 text-slate-300 hover:bg-slate-800/80'
                    : 'bg-slate-950/40 border-slate-800/60 text-slate-500 hover:bg-slate-900/40'
                }`}
              >
                <div className="flex items-center justify-between w-full mb-1">
                  <span className={`text-[10px] font-mono font-bold ${isCurrent ? 'text-indigo-400' : isCompleted ? 'text-emerald-400' : 'text-slate-500'}`}>
                    STEP 0{step.step}
                  </span>
                  {isCompleted ? (
                    <CheckCircle2 className="w-3.5 h-3.5 text-emerald-400" />
                  ) : isCurrent ? (
                    <div className="w-2 h-2 rounded-full bg-indigo-400 animate-ping" />
                  ) : (
                    <Clock className="w-3.5 h-3.5 text-slate-600" />
                  )}
                </div>
                <span className="text-xs font-semibold leading-tight line-clamp-1 text-slate-200">
                  {step.title}
                </span>
              </button>
            );
          })}
        </div>
      </div>

      {/* Main Execution View: Split into Details & Live Telemetry Inspector */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Left: Step Explanation & 3D State Vector (7 cols) */}
        <div className="lg:col-span-7 space-y-6">
          {/* Step Detail Card */}
          <div className="rounded-2xl border border-slate-800 bg-slate-900/50 p-6 space-y-5">
            <div className="flex items-center justify-between">
              <span className="px-2.5 py-1 rounded bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 text-xs font-mono font-semibold uppercase">
                {currentStep.category}
              </span>
              <span className="text-xs text-slate-400 font-mono">
                Step {currentStep.step} of 7
              </span>
            </div>

            <div>
              <h2 className="text-xl font-bold text-slate-100">{currentStep.title}</h2>
              <p className="text-sm text-slate-300 mt-2 leading-relaxed">
                {currentStep.description}
              </p>
            </div>

            <div className="p-3.5 rounded-xl bg-slate-950/80 border border-slate-800 flex items-center justify-between font-mono text-xs">
              <div className="flex items-center gap-2">
                <span className={`px-2 py-0.5 rounded font-bold ${
                  currentStep.httpMethod === 'POST' ? 'bg-cyan-950 text-cyan-400 border border-cyan-800/60' : 'bg-emerald-950 text-emerald-400 border border-emerald-800/60'
                }`}>
                  {currentStep.httpMethod}
                </span>
                <span className="text-slate-300">{currentStep.endpoint}</span>
              </div>
              <span className="text-slate-500 text-[11px]">HTTP REST</span>
            </div>

            {/* Highlights depending on step */}
            {currentStep.step === 4 && (
              <div className="p-3.5 rounded-xl bg-rose-950/20 border border-rose-800/40 flex items-start gap-3">
                <AlertTriangle className="w-5 h-5 text-rose-400 shrink-0 mt-0.5" />
                <div className="text-xs">
                  <span className="font-semibold text-rose-300">Crucial Defended Gate:</span>
                  <p className="text-slate-400 mt-0.5">
                    A concurrent second capture request during this window will be intercepted by PayBridge's in-process mutex lock, returning <strong className="text-slate-200">HTTP 409 Conflict</strong>.
                  </p>
                </div>
              </div>
            )}

            {currentStep.step === 6 && (
              <div className="p-3.5 rounded-xl bg-emerald-950/20 border border-emerald-800/40 flex items-start gap-3">
                <CheckCircle2 className="w-5 h-5 text-emerald-400 shrink-0 mt-0.5" />
                <div className="text-xs">
                  <span className="font-semibold text-emerald-300">Zero Payment Mutations:</span>
                  <p className="text-slate-400 mt-0.5">
                    Reconciliation executed strictly via reader-only queries. No money moved, no duplicate authorization created.
                  </p>
                </div>
              </div>
            )}

            {currentStep.step === 7 && (
              <div className="p-4 rounded-xl bg-indigo-950/30 border border-indigo-800/40 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <Award className="w-6 h-6 text-indigo-400" />
                  <div>
                    <h3 className="text-xs font-bold text-slate-100">Truth Passport Ready</h3>
                    <p className="text-[11px] text-slate-400">Cryptographically verifiable evidence bundle issued.</p>
                  </div>
                </div>
                <button
                  onClick={() => onNavigate('passport')}
                  className="px-3 py-1.5 rounded-lg bg-indigo-600 hover:bg-indigo-500 text-white text-xs font-semibold transition-all flex items-center gap-1.5 shadow-sm"
                >
                  <span>View Passport</span>
                  <ArrowRight className="w-3.5 h-3.5" />
                </button>
              </div>
            )}
          </div>

          {/* Current 3D State Vector Card */}
          <div className="rounded-2xl border border-slate-800 bg-slate-900/50 p-6 space-y-4">
            <h3 className="text-xs font-bold text-slate-400 uppercase tracking-wider font-mono">
              3D State Vector (After Step {currentStep.step})
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div className="p-3 rounded-xl bg-slate-950/60 border border-slate-800">
                <span className="text-[11px] text-slate-500 block mb-1">Provider State</span>
                <StatusPill type="provider" status={currentStep.stateVectorAfter.provider_state} />
              </div>
              <div className="p-3 rounded-xl bg-slate-950/60 border border-slate-800">
                <span className="text-[11px] text-slate-500 block mb-1">Application State</span>
                <StatusPill type="app" status={currentStep.stateVectorAfter.application_state} />
              </div>
              <div className="p-3 rounded-xl bg-slate-950/60 border border-slate-800">
                <span className="text-[11px] text-slate-500 block mb-1">Case Status</span>
                <StatusPill type="case" status={currentStep.stateVectorAfter.case_status} />
              </div>
            </div>
          </div>
        </div>

        {/* Right: Telemetry & Payload Inspector (5 cols) */}
        <div className="lg:col-span-5 rounded-2xl border border-slate-800 bg-slate-950/90 flex flex-col h-[520px] overflow-hidden shadow-xl">
          {/* Header & Tabs */}
          <div className="flex items-center justify-between border-b border-slate-800 px-4 py-3 bg-slate-900/70">
            <div className="flex items-center gap-2">
              <Code className="w-4 h-4 text-cyan-400" />
              <span className="text-xs font-bold font-mono text-slate-300">TELEMETRY INSPECTOR</span>
            </div>
            <div className="flex items-center gap-1 bg-slate-950 p-1 rounded-lg border border-slate-800">
              <button
                onClick={() => setActiveTab('response')}
                className={`px-2.5 py-1 rounded text-[11px] font-mono transition-all ${
                  activeTab === 'response' ? 'bg-indigo-600 text-white font-semibold' : 'text-slate-400 hover:text-slate-200'
                }`}
              >
                Response
              </button>
              {currentStep.requestPayload && (
                <button
                  onClick={() => setActiveTab('request')}
                  className={`px-2.5 py-1 rounded text-[11px] font-mono transition-all ${
                    activeTab === 'request' ? 'bg-indigo-600 text-white font-semibold' : 'text-slate-400 hover:text-slate-200'
                  }`}
                >
                  Request
                </button>
              )}
            </div>
          </div>

          {/* JSON Body */}
          <div className="flex-1 p-4 overflow-y-auto font-mono text-xs text-slate-300 bg-slate-950/60 leading-relaxed selection:bg-indigo-500 selection:text-white">
            <pre className="whitespace-pre-wrap">
              {JSON.stringify(
                activeTab === 'response' ? currentStep.responsePayload : currentStep.requestPayload,
                null,
                2
              )}
            </pre>
          </div>

          {/* Footer Metadata */}
          <div className="border-t border-slate-800 px-4 py-2.5 bg-slate-900/40 flex items-center justify-between text-[11px] text-slate-400 font-mono">
            <span>Status: 200 OK</span>
            <span>Duration: 18ms</span>
          </div>
        </div>
      </div>
    </div>
  );
};
