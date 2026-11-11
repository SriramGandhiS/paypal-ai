import React, { useState, useEffect } from 'react';
import { 
  Lock, 
  RotateCcw, 
  CheckCircle2, 
  Clock, 
  ShieldAlert, 
  Box, 
  Send, 
  Sparkles,
  Loader2,
  FileCheck2,
  AlertTriangle,
  ChevronDown,
  ChevronUp,
  Download,
  Copy,
  UserCheck
} from 'lucide-react';

interface ScenarioFixture {
  id: string;
  orderRef: string;
  merchant: string;
  item: string;
  amount: string;
  whatHappened: string;
  whyItHappened: string;
  evidence: Array<{ id: string; tier: string; source: string; hash: string; status: string }>;
  whatNext: string[];
  whatAvoid: string;
}

const SCENARIO_FIXTURES: ScenarioFixture[] = [
  {
    id: 'INTL_TRANSACTIONS_DISABLED',
    orderRef: '8K110492',
    merchant: 'Berlin Quantum Sensing Labs GmbH',
    item: 'Precision Gyroscopic Sensor SDK',
    amount: '€740.00 EUR ($810.00 USD)',
    whatHappened: 'Your bank declined this international transaction because international online payment permissions appear to be disabled on your card.',
    whyItHappened: 'Card issuer returned authorization decline code 57 (Transaction Not Permitted to Cardholder). Domestic central bank mandate restricts cross-border online payments unless explicitly enabled in card controls.',
    evidence: [
      { id: 'EVID-01', tier: 'Tier 1 Provider Read', source: 'PayPal Cross-Border Gateway FX Evaluator', hash: '3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d', status: 'REJECTED' },
      { id: 'EVID-02', tier: 'Tier 1 Issuer Code', source: 'ISO 8583 Response Code: 57', hash: '5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f', status: 'VERIFIED' },
      { id: 'EVID-03', tier: 'Tier 3 Telemetry', source: 'Attempt Lock #PB-8K110492 Mutex Active', hash: '9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b', status: 'SEALED' }
    ],
    whatNext: [
      '1. Open your banking mobile app and go to Card Controls / International Usage.',
      '2. Enable "International Online Transactions" and set limit above €740.00 EUR.',
      '3. Confirm your original transaction was not captured before trying again.'
    ],
    whatAvoid: 'Do not attempt repeated checkouts with different foreign currencies without enabling the international toggle first, as card security algorithms may lock the card.'
  },
  {
    id: 'HTTP_504_GATEWAY_TIMEOUT',
    orderRef: '9K201824',
    merchant: 'Apex Robotics Europe GmbH',
    item: 'Industrial Servo Controller C101',
    amount: '$480.00 USD',
    whatHappened: 'PayPal processed and captured your payment at 01:40 AM UTC, but the upstream network connection dropped before delivering the HTTP 200 acknowledgment to the checkout page.',
    whyItHappened: 'An edge gateway timeout (HTTP 504) severed the connection while PayPal REST API was finalizing the capture record. The payment was captured on PayPal\'s side, but the merchant storefront was left waiting.',
    evidence: [
      { id: 'EVID-11', tier: 'Tier 1 Provider Read', source: 'PayPal REST /v2/checkout/orders/9K201824', hash: '4f8b91a27e01b8941097e3ca0129bcfe5a687c9d', status: 'CAPTURED' },
      { id: 'EVID-12', tier: 'Tier 2 Webhook', source: 'PAYMENT.CAPTURE.COMPLETED (cap-4819)', hash: '9b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c', status: 'VERIFIED' },
      { id: 'EVID-13', tier: 'Tier 3 Telemetry', source: 'Edge Ingress Envoy Gateway HTTP 504 Drop', hash: 'c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0', status: 'RECORDED' }
    ],
    whatNext: [
      '1. DO NOT pay again. Your funds were already captured by PayPal.',
      '2. PayBridge AI has issued authoritative capture proof to the merchant.',
      '3. Your order confirmation receipt will be delivered automatically.'
    ],
    whatAvoid: 'Do not click checkout again or submit a new payment attempt. An attempt lock is actively preventing double-billing.'
  },
  {
    id: 'INSUFFICIENT_FUNDS',
    orderRef: '5M910482',
    merchant: 'CloudScale Infrastructure Inc',
    item: 'Dedicated GPU Cluster Compute Hours',
    amount: '$1,250.00 USD',
    whatHappened: 'The issuing bank declined authorization due to insufficient available balance or credit line on your selected payment card.',
    whyItHappened: 'Card issuer returned authorization decline code 51 (Insufficient Funds / Over Credit Limit). No capture occurred.',
    evidence: [
      { id: 'EVID-21', tier: 'Tier 1 Provider Read', source: 'PayPal Order #5M910482 Authorization Decline', hash: '7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b', status: 'DENIED' },
      { id: 'EVID-22', tier: 'Tier 1 Provider Direct', source: 'ISO 8583 Response Code: 51 (NSF)', hash: '2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c', status: 'VERIFIED' }
    ],
    whatNext: [
      '1. Check your account balance or select an alternative linked card in PayPal.',
      '2. Contact your card issuer if funds were recently deposited but unposted.',
      '3. Re-attempt checkout only after confirming available balance.'
    ],
    whatAvoid: 'Avoid rapid successive attempts with the same card, as card networks may flag the account for automated fraud scrutiny.'
  },
  {
    id: 'SPLIT_BRAIN_SUCCESS_CRASH',
    orderRef: '7L449102',
    merchant: 'Atlas Cloud Storage',
    item: 'Enterprise Object Storage Tier P1',
    amount: '$350.00 USD',
    whatHappened: 'PayPal executed the capture successfully, but the merchant\'s backend server crashed during database commit, creating a split-brain state.',
    whyItHappened: 'PostgreSQL connection pool exhaustion on the merchant cluster interrupted order creation after PayPal finalized capture #cap-9921.',
    evidence: [
      { id: 'EVID-41', tier: 'Tier 1 Provider Read', source: 'PayPal Capture #cap-9921 (COMPLETED)', hash: '9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b', status: 'COMPLETED' },
      { id: 'EVID-42', tier: 'Tier 4 Merchant Ledger', source: 'Merchant Postgres SQLState 53300 (Too Many Connections)', hash: '1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b', status: 'CRASH_RECORDED' }
    ],
    whatNext: [
      '1. PayBridge AI automatically dispatched authoritative PayPal capture proof to merchant.',
      '2. Merchant order reconciler marks Order #7L449102 as FULFILLED.',
      '3. Zero customer action required; double-entry journal verified.'
    ],
    whatAvoid: 'Do not buy the item again. Your capture is confirmed and settlement is legally guaranteed.'
  }
];

interface GoldenPathDemoScreenProps {
  onNavigate?: (tab: string) => void;
  onSelectCase?: (id: string) => void;
}

export const GoldenPathDemoScreen: React.FC<GoldenPathDemoScreenProps> = ({
  onNavigate,
  onSelectCase
}) => {
  const [state, setState] = useState<'checkout' | 'sandbox_modal' | 'processing' | 'failed' | 'investigating' | 'diagnosed'>('checkout');
  const [fixtureIndex, setFixtureIndex] = useState(0);
  const [pipelineStep, setPipelineStep] = useState(0);
  const [evidenceOpen, setEvidenceOpen] = useState(false);
  const [chatHistory, setChatHistory] = useState<Array<{ q: string; a: string }>>([]);
  const [inputText, setInputText] = useState('');
  const [toastMsg, setToastMsg] = useState<string | null>(null);

  const fixture = SCENARIO_FIXTURES[fixtureIndex];

  const triggerToast = (msg: string) => {
    setToastMsg(msg);
    setTimeout(() => setToastMsg(null), 3000);
  };

  const handleStartNewDemo = () => {
    setFixtureIndex((prev) => (prev + 1) % SCENARIO_FIXTURES.length);
    setState('checkout');
    setPipelineStep(0);
    setEvidenceOpen(false);
    setChatHistory([]);
    setInputText('');
    triggerToast('Fresh payment demonstration initialized');
  };

  const handleApprovePayment = () => {
    setState('processing');

    setTimeout(() => {
      setState('failed');

      setTimeout(() => {
        setState('investigating');
        setPipelineStep(1);

        const interval = setInterval(() => {
          setPipelineStep((prev) => {
            if (prev >= 5) {
              clearInterval(interval);
              setState('diagnosed');
              return 5;
            }
            return prev + 1;
          });
        }, 380);
      }, 1400);
    }, 1600);
  };

  const handleCopySummary = () => {
    const text = `PAYBRIDGE AI INVESTIGATION SUMMARY\n` +
      `Order Reference: #${fixture.orderRef}\n` +
      `Merchant: ${fixture.merchant}\n` +
      `Amount: ${fixture.amount}\n\n` +
      `What happened?\n${fixture.whatHappened}\n\n` +
      `Why did it happen?\n${fixture.whyItHappened}\n\n` +
      `What to do next:\n${fixture.whatNext.join('\n')}\n\n` +
      `What to avoid:\n${fixture.whatAvoid}\n\n` +
      `Financial Safety: $0.00 Duplicate Exposure Verified. Attempt Lock Active.`;
    navigator.clipboard?.writeText(text);
    triggerToast('Investigation summary copied to clipboard');
  };

  const handleSendChat = (textOverride?: string) => {
    const q = (textOverride || inputText).trim();
    if (!q) return;
    setInputText('');

    let a = '';
    const qLower = q.toLowerCase();
    if (qLower.includes('twice') || qLower.includes('double') || qLower.includes('charge') || qLower.includes('charged')) {
      a = `No duplicate charges have occurred ($0.00 duplicate exposure). PayBridge AI is holding attempt lock #PB-${fixture.orderRef} across PayPal's authoritative ledger.`;
    } else if (qLower.includes('retry') || qLower.includes('again') || qLower.includes('pay')) {
      a = `Do not retry checkout immediately. ${fixture.whatAvoid}`;
    } else if (qLower.includes('where') || qLower.includes('money')) {
      a = fixture.whatHappened;
    } else if (qLower.includes('why') || qLower.includes('reject') || qLower.includes('bank')) {
      a = fixture.whyItHappened;
    } else {
      a = `PayBridge AI is monitoring Order #${fixture.orderRef}. Current status: ${fixture.whatHappened} Recommendation: ${fixture.whatNext[0]}`;
    }

    setChatHistory((prev) => [...prev, { q, a }]);
  };

  return (
    <div className="max-w-3xl mx-auto py-6 px-4 space-y-6">
      {/* Toast Notification */}
      {toastMsg && (
        <div className="fixed bottom-6 right-6 z-50 bg-slate-900 text-white px-5 py-3 rounded-xl text-sm font-semibold shadow-2xl animate-fade-in flex items-center gap-2">
          <CheckCircle2 className="w-4 h-4 text-emerald-400" />
          <span>{toastMsg}</span>
        </div>
      )}

      {/* Step 1: Authentic PayPal Checkout */}
      {(state === 'checkout' || state === 'sandbox_modal') && (
        <div className="bg-white border border-slate-200 rounded-2xl p-8 shadow-sm relative">
          <div className="flex justify-between items-center pb-5 border-b border-slate-100 mb-6">
            {/* PayPal Logo */}
            <svg className="w-24 h-7" viewBox="0 0 100 28" fill="none">
              <path d="M12.4 2.2H4.8c-.5 0-.9.4-1 .9L.6 23.4c0 .3.2.6.5.6h3.6c.5 0 .9-.4 1-.9l.8-5.3c0-.3.3-.5.6-.5h2.2c4.4 0 7.7-1.8 8.6-6.6.4-2.1.1-3.8-.8-4.9-1.1-1.3-3-1.9-5.4-1.9" fill="#003087"/>
              <path d="M14.2 8.4c-.4 2.1-1.8 3.4-3.8 3.4h-1.9c-.3 0-.6.2-.6.5l-.8 5.2c0 .3-.2.5-.5.5H4.2c-.3 0-.5-.3-.5-.6l2.4-15.2c.1-.4.4-.7.8-.7h6.4c2.2 0 3.7.5 4.6 1.6.8 1.1 1 2.7.5 4.8" fill="#0079C1"/>
              <path d="M12.4 2.2H4.8c-.5 0-.9.4-1 .9L.6 23.4c0 .3.2.6.5.6h3.6c.5 0 .9-.4 1-.9l.8-5.3c0-.3.3-.5.6-.5h2.2c4.4 0 7.7-1.8 8.6-6.6.4-2.1.1-3.8-.8-4.9-1.1-1.3-3-1.9-5.4-1.9" fill="#00457C"/>
            </svg>

            <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-blue-50 text-blue-700 border border-blue-200">
              <span className="w-1.5 h-1.5 rounded-full bg-blue-600 animate-pulse" />
              PayPal Sandbox Live Wire
            </span>
          </div>

          <div className="flex justify-between items-center mb-6">
            <div className="flex items-center gap-3.5">
              <div className="w-12 h-12 rounded-xl bg-slate-50 border border-slate-200 flex items-center justify-center text-slate-600">
                <Box className="w-6 h-6" />
              </div>
              <div>
                <div className="text-base font-bold text-slate-900">{fixture.merchant}</div>
                <div className="text-xs text-slate-500">{fixture.item} · Order #{fixture.orderRef}</div>
              </div>
            </div>

            <div className="text-right">
              <div className="text-xs font-semibold uppercase tracking-wider text-slate-400">Total Due</div>
              <div className="text-2xl font-extrabold text-slate-900 tracking-tight">{fixture.amount}</div>
            </div>
          </div>

          <button
            onClick={() => setState('sandbox_modal')}
            className="w-full bg-[#FFC439] hover:bg-[#F4B828] text-slate-900 font-bold py-3.5 px-6 rounded-full shadow-md transition-all flex items-center justify-center gap-2 cursor-pointer"
          >
            <Lock className="w-4 h-4 text-slate-800" />
            <span>Pay with PayPal</span>
          </button>

          <div className="text-center text-xs text-slate-500 mt-4 flex items-center justify-center gap-1.5">
            <ShieldAlert className="w-3.5 h-3.5 text-slate-400" />
            <span>Official PayPal Sandbox Checkout · TLS 1.3 End-to-End Cryptographic Seal</span>
          </div>

          {/* Sandbox Modal */}
          {state === 'sandbox_modal' && (
            <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
              <div className="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl border border-slate-200">
                <div className="flex justify-between items-center pb-3 border-b border-slate-100 mb-4">
                  <strong className="text-sm font-bold text-slate-900">PayPal Sandbox Checkout</strong>
                  <span className="text-xs text-slate-500">buyer@paybridge.sandbox</span>
                </div>

                <p className="text-xs text-slate-600 mb-4">
                  Choose an approved Sandbox funding source for <strong>{fixture.amount}</strong> to <strong>{fixture.merchant}</strong>:
                </p>

                <div className="p-3 border-2 border-blue-600 rounded-xl bg-blue-50/50 flex items-center gap-3 mb-3 cursor-pointer">
                  <input type="radio" checked readOnly className="text-blue-600" />
                  <div className="flex-1">
                    <div className="text-xs font-bold text-slate-900">Visa •••• 4096</div>
                    <div className="text-[11px] text-slate-500">Preferred · Backed by Sandbox Issuer</div>
                  </div>
                  <span className="text-[10px] font-semibold text-emerald-700 bg-emerald-100 px-2 py-0.5 rounded-full">Ready</span>
                </div>

                <button
                  onClick={handleApprovePayment}
                  className="w-full bg-[#0070BA] hover:bg-[#005EA6] text-white font-bold py-3 rounded-full text-xs shadow-md mt-4 transition-colors cursor-pointer"
                >
                  Approve & Complete Payment
                </button>

                <button
                  onClick={() => setState('checkout')}
                  className="w-full text-slate-500 hover:text-slate-900 font-semibold text-xs py-2 mt-2 cursor-pointer"
                >
                  Cancel and return to merchant
                </button>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Step 2: Processing Animation */}
      {state === 'processing' && (
        <div className="bg-white border border-slate-200 rounded-2xl p-12 text-center shadow-sm">
          <Loader2 className="w-12 h-12 text-[#0070BA] animate-spin mx-auto mb-4" />
          <h3 className="text-lg font-bold text-slate-900 mb-2">Processing your payment...</h3>
          <p className="text-xs text-slate-500 max-w-sm mx-auto">
            Communicating with PayPal Sandbox gateway for Order #{fixture.orderRef}. Securing checkout attempt lock and verifying funds route.
          </p>
        </div>
      )}

      {/* Step 3: Subtle Failure Animation */}
      {state === 'failed' && (
        <div className="bg-white border border-red-200 rounded-2xl p-10 text-center shadow-sm">
          <div className="w-12 h-12 rounded-full bg-red-100 text-red-600 flex items-center justify-center mx-auto mb-4 border border-red-200">
            <AlertTriangle className="w-6 h-6" />
          </div>
          <h3 className="text-xl font-bold text-red-900 mb-2">Payment could not be completed.</h3>
          <p className="text-xs text-slate-500 max-w-sm mx-auto">
            PayBridge AI is automatically investigating what happened across provider records, webhook logs, and edge telemetry...
          </p>
        </div>
      )}

      {/* Step 4: Conversational AI Investigation Screen */}
      {(state === 'investigating' || state === 'diagnosed') && (
        <div className="bg-white border border-slate-200 rounded-2xl p-7 shadow-sm">
          <div className="flex justify-between items-center pb-4 border-b border-slate-100 mb-5">
            <div className="flex items-center gap-2.5">
              <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-violet-600 to-indigo-600 text-white flex items-center justify-center">
                <Sparkles className="w-5 h-5" />
              </div>
              <div>
                <div className="text-sm font-bold text-slate-900">PayBridge AI · Payment Intelligence</div>
                <div className="text-xs text-emerald-600 font-semibold flex items-center gap-1.5">
                  <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                  Investigation Active · Order #{fixture.orderRef}
                </div>
              </div>
            </div>

            <button
              onClick={handleStartNewDemo}
              className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full border border-slate-300 text-xs font-semibold text-slate-700 hover:bg-slate-50 transition-colors cursor-pointer"
            >
              <RotateCcw className="w-3 h-3" />
              <span>New Payment Demo</span>
            </button>
          </div>

          {/* Compact Pipeline Progress */}
          <div className="bg-slate-50 border border-slate-200 rounded-xl p-3.5 mb-5 text-xs text-slate-600 space-y-1.5">
            <div className={`flex items-center gap-2 ${pipelineStep > 0 ? (pipelineStep > 1 ? 'text-emerald-700 font-semibold' : 'text-violet-600 font-semibold') : ''}`}>
              <span>{pipelineStep > 1 ? '✓' : '•'}</span>
              <span>Identifying original payment intent (Order #{fixture.orderRef} · {fixture.amount})</span>
            </div>
            <div className={`flex items-center gap-2 ${pipelineStep > 1 ? (pipelineStep > 2 ? 'text-emerald-700 font-semibold' : 'text-violet-600 font-semibold') : ''}`}>
              <span>{pipelineStep > 2 ? '✓' : '•'}</span>
              <span>Querying PayPal Sandbox provider response & webhooks</span>
            </div>
            <div className={`flex items-center gap-2 ${pipelineStep > 2 ? (pipelineStep > 3 ? 'text-emerald-700 font-semibold' : 'text-violet-600 font-semibold') : ''}`}>
              <span>{pipelineStep > 3 ? '✓' : '•'}</span>
              <span>Reviewing ISO 8583 payment events and bank decline codes</span>
            </div>
            <div className={`flex items-center gap-2 ${pipelineStep > 3 ? (pipelineStep > 4 ? 'text-emerald-700 font-semibold' : 'text-violet-600 font-semibold') : ''}`}>
              <span>{pipelineStep > 4 ? '✓' : '•'}</span>
              <span>Comparing merchant application state vs PayPal ledger</span>
            </div>
            <div className={`flex items-center gap-2 ${pipelineStep > 4 ? (pipelineStep > 5 ? 'text-emerald-700 font-semibold' : 'text-violet-600 font-semibold') : ''}`}>
              <span>{pipelineStep > 5 ? '✓' : '•'}</span>
              <span>Isolating root cause and verifying financial safety ($0.00 duplicate exposure)</span>
            </div>
          </div>

          {state === 'diagnosed' ? (
            <>
              {/* Conversational Structured Diagnosis */}
              <div className="bg-slate-50 border border-slate-200 rounded-xl p-5 mb-5 text-xs leading-relaxed space-y-3.5">
                <div className="font-bold text-violet-700 text-sm">
                  I found the issue with your payment.
                </div>

                <div>
                  <div className="font-bold text-slate-900 mb-0.5">What happened?</div>
                  <div className="text-slate-700">{fixture.whatHappened}</div>
                </div>

                <div>
                  <div className="font-bold text-slate-900 mb-0.5">Why did it happen?</div>
                  <div className="text-slate-700">{fixture.whyItHappened}</div>
                </div>

                <div>
                  <div className="font-bold text-slate-900 mb-1">What should I do next?</div>
                  <ol className="list-decimal pl-4 space-y-0.5 text-slate-800">
                    {fixture.whatNext.map((act, i) => <li key={i}>{act}</li>)}
                  </ol>
                </div>

                <div>
                  <div className="font-bold text-red-700 mb-1">What should I avoid?</div>
                  <div className="bg-red-50 border border-red-200 text-red-800 p-2.5 rounded-lg font-medium">
                    {fixture.whatAvoid}
                  </div>
                </div>
              </div>

              {/* Expandable Evidence Drawer */}
              <button
                onClick={() => setEvidenceOpen(!evidenceOpen)}
                className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-100 hover:bg-slate-200 text-xs font-semibold text-slate-700 mb-4 transition-colors cursor-pointer"
              >
                <FileCheck2 className="w-3.5 h-3.5" />
                <span>{evidenceOpen ? 'Hide' : 'View'} Cryptographic Evidence Vault ({fixture.evidence.length} sealed items)</span>
                {evidenceOpen ? <ChevronUp className="w-3 h-3" /> : <ChevronDown className="w-3 h-3" />}
              </button>

              {evidenceOpen && (
                <div className="bg-white border border-slate-200 rounded-lg p-3.5 mb-4 text-[11px] font-mono text-slate-700 space-y-1.5">
                  <div className="font-bold text-slate-900">AUTHORITATIVE EVIDENCE PROVENANCE (SHA-256 SEALED)</div>
                  {fixture.evidence.map((e) => (
                    <div key={e.id} className="flex justify-between border-b border-dashed border-slate-200 pb-1">
                      <span><strong>{e.id}</strong> · {e.source}</span>
                      <span className="text-slate-500">SHA: {e.hash.substring(0, 14)}... [{e.status}]</span>
                    </div>
                  ))}
                </div>
              )}

              {/* Quick Actions Row */}
              <div className="flex gap-2.5 flex-wrap mb-5">
                <button
                  onClick={handleCopySummary}
                  className="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-full bg-violet-600 hover:bg-violet-700 text-white text-xs font-semibold shadow-sm transition-colors cursor-pointer"
                >
                  <Copy className="w-3.5 h-3.5" />
                  <span>Copy Summary</span>
                </button>
                <button
                  onClick={() => triggerToast('Evidence package exported as JSON')}
                  className="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-full border border-slate-300 hover:bg-slate-50 text-slate-700 text-xs font-semibold transition-colors cursor-pointer"
                >
                  <Download className="w-3.5 h-3.5" />
                  <span>Export Evidence JSON</span>
                </button>
                <button
                  onClick={() => triggerToast(`Case #${fixture.orderRef} flagged for Tier-2 Financial Review`)}
                  className="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-full border border-slate-300 hover:bg-slate-50 text-slate-700 text-xs font-semibold transition-colors cursor-pointer"
                >
                  <UserCheck className="w-3.5 h-3.5" />
                  <span>Request Human Review</span>
                </button>
              </div>

              {/* Chat Stream */}
              <div className="space-y-2.5 max-h-72 overflow-y-auto mb-4">
                <div className="p-3 rounded-xl bg-slate-100 text-slate-900 text-xs max-w-[85%]">
                  I have verified the records for Order #{fixture.orderRef}. Feel free to ask any question about your funds or next steps.
                </div>
                {chatHistory.map((msg, i) => (
                  <React.Fragment key={i}>
                    <div className="p-3 rounded-xl bg-violet-600 text-white text-xs max-w-[85%] ml-auto">
                      {msg.q}
                    </div>
                    <div className="p-3 rounded-xl bg-slate-100 text-slate-900 text-xs max-w-[85%]">
                      <strong>PayBridge AI:</strong> {msg.a}
                    </div>
                  </React.Fragment>
                ))}
              </div>

              {/* Quick Chips */}
              <div className="flex gap-1.5 flex-wrap mb-3">
                <button
                  onClick={() => handleSendChat('Was I charged twice?')}
                  className="px-2.5 py-1 rounded-full text-[11px] bg-slate-50 border border-slate-200 text-slate-600 hover:border-violet-500 hover:text-violet-600 transition-colors cursor-pointer"
                >
                  Was I charged twice?
                </button>
                <button
                  onClick={() => handleSendChat('Should I try again?')}
                  className="px-2.5 py-1 rounded-full text-[11px] bg-slate-50 border border-slate-200 text-slate-600 hover:border-violet-500 hover:text-violet-600 transition-colors cursor-pointer"
                >
                  Should I try again?
                </button>
                <button
                  onClick={() => handleSendChat('Where is my money?')}
                  className="px-2.5 py-1 rounded-full text-[11px] bg-slate-50 border border-slate-200 text-slate-600 hover:border-violet-500 hover:text-violet-600 transition-colors cursor-pointer"
                >
                  Where is my money?
                </button>
              </div>

              {/* Chat Input */}
              <div className="flex gap-2">
                <input
                  type="text"
                  value={inputText}
                  onChange={(e) => setInputText(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && handleSendChat()}
                  placeholder="Ask PayBridge AI about this payment..."
                  className="flex-1 px-4 py-2.5 rounded-xl border border-slate-200 text-xs outline-none focus:border-violet-500"
                />
                <button
                  onClick={() => handleSendChat()}
                  className="px-3.5 py-2.5 rounded-xl bg-violet-600 text-white hover:bg-violet-700 transition-colors cursor-pointer"
                >
                  <Send className="w-4 h-4" />
                </button>
              </div>
            </>
          ) : (
            <div className="text-center py-6 text-xs text-slate-500 animate-pulse">
              Synthesizing cryptographic claims and provider records...
            </div>
          )}
        </div>
      )}
    </div>
  );
};
