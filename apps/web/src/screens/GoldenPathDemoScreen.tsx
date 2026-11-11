import React, { useState } from 'react';
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
  FileCheck2
} from 'lucide-react';
import { api } from '../api/client';

interface GoldenPathDemoScreenProps {
  onNavigate?: (tab: string) => void;
  onSelectCase?: (caseId: string) => void;
}

export const GoldenPathDemoScreen: React.FC<GoldenPathDemoScreenProps> = ({
  onNavigate,
  onSelectCase
}) => {
  const [demoState, setDemoState] = useState<'idle' | 'initiating' | 'timed_out' | 'reconciling' | 'reconciled'>('idle');
  const [chatHistory, setChatHistory] = useState<Array<{ q: string; a: string }>>([]);
  const [inputText, setInputText] = useState('');

  const handleStartCheckout = () => {
    if (demoState !== 'idle') return;
    setDemoState('initiating');

    setTimeout(() => {
      setDemoState('timed_out');

      setTimeout(() => {
        setDemoState('reconciling');

        setTimeout(() => {
          setDemoState('reconciled');
        }, 2600);
      }, 1800);
    }, 1400);
  };

  const handleReset = () => {
    setDemoState('idle');
    setChatHistory([]);
    setInputText('');
  };

  const handleSendQuery = (textOverride?: string) => {
    const q = (textOverride || inputText).trim();
    if (!q) return;
    setInputText('');

    let answer = '';
    const qLower = q.toLowerCase();
    if (qLower.includes('retry') || qLower.includes('again')) {
      answer = "Do not pay again yet. A duplicate checkout would capture an unnecessary second $480.00 charge. We are holding the checkout mutex lock while inspecting PayPal's authoritative ledger.";
    } else if (qLower.includes('charge') || qLower.includes('double') || qLower.includes('charged')) {
      answer = "No duplicate charge has occurred ($0.00 duplicate exposure). The original payment intent is verified with PayPal REST API order #9K201824 and capture #cap-4819.";
    } else if (qLower.includes('where') || qLower.includes('money') || qLower.includes('status')) {
      answer = "Your $480.00 USD was successfully captured by PayPal Sandbox at 01:40 AM. The HTTP 504 was purely an edge network timeout. We have confirmed the funds and will complete settlement without extra fees.";
    } else if (qLower.includes('what happened') || qLower.includes('504') || qLower.includes('timeout')) {
      answer = "PayPal processed the transaction, but the upstream connection dropped before sending the HTTP 200 acknowledgment. PayBridge AI automatically intercepted the timeout to protect your funds.";
    } else {
      answer = "PayBridge AI is monitoring session #PB-9482. The original payment of $480.00 USD is verified against PayPal's ledger. Attempt locks prevent duplicate billing.";
    }

    setChatHistory(prev => [...prev, { q, a: answer }]);
  };

  return (
    <div className="max-w-5xl mx-auto py-4 px-2 space-y-7 animate-fade-in">
      {/* Section A: Header */}
      <div>
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white border border-slate-200 text-xs font-medium text-slate-900 shadow-sm mb-3">
          <span className="w-2.5 h-2.5 rounded-full border-2 border-violet-500 shadow-sm shadow-violet-500/50" />
          <span>Golden Path Demo</span>
        </div>
        <h1 className="text-3xl font-bold tracking-tight text-slate-900">
          Cross-Border 504 Recovery <span className="text-violet-600 bg-gradient-to-r from-violet-600 to-purple-600 bg-clip-text text-transparent">in Action.</span>
        </h1>
        <p className="text-sm text-slate-500 mt-1.5">
          A live simulation of a PayPal checkout timeout and autonomous recovery by PayBridge AI.
        </p>
      </div>

      {/* Section B & C: Side-by-Side Cards (Matches Reference Image) */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        
        {/* Left Card: PayPal Checkout Experience */}
        <div className="lg:col-span-5 bg-white border border-slate-200 rounded-2xl p-6 shadow-sm">
          {/* PayPal Logo */}
          <div className="pb-4 border-b border-slate-100 mb-5">
            <svg className="w-24 h-7" viewBox="0 0 100 28" fill="none">
              <path d="M12.4 2.2H4.8c-.5 0-.9.4-1 .9L.6 23.4c0 .3.2.6.5.6h3.6c.5 0 .9-.4 1-.9l.8-5.3c0-.3.3-.5.6-.5h2.2c4.4 0 7.7-1.8 8.6-6.6.4-2.1.1-3.8-.8-4.9-1.1-1.3-3-1.9-5.4-1.9" fill="#003087"/>
              <path d="M14.2 8.4c-.4 2.1-1.8 3.4-3.8 3.4h-1.9c-.3 0-.6.2-.6.5l-.8 5.2c0 .3-.2.5-.5.5H4.2c-.3 0-.5-.3-.5-.6l2.4-15.2c.1-.4.4-.7.8-.7h6.4c2.2 0 3.7.5 4.6 1.6.8 1.1 1 2.7.5 4.8" fill="#0079C1"/>
              <path d="M12.4 2.2H4.8c-.5 0-.9.4-1 .9L.6 23.4c0 .3.2.6.5.6h3.6c.5 0 .9-.4 1-.9l.8-5.3c0-.3.3-.5.6-.5h2.2c4.4 0 7.7-1.8 8.6-6.6.4-2.1.1-3.8-.8-4.9-1.1-1.3-3-1.9-5.4-1.9" fill="#00457C"/>
              <path d="M28.4 9.8h-3.4c-.3 0-.6.2-.6.5l-.2 1.1c-.5-.7-1.5-1.3-2.9-1.3-3.1 0-5.8 2.4-6.3 5.7-.3 1.7 0 3.3.9 4.5 1 1.2 2.5 1.8 4.2 1.8 2.5 0 3.9-1.6 3.9-1.6l-.2 1.1c0 .3.2.6.5.6h3.1c.4 0 .7-.3.7-.7l1.9-12.1c0-.4-.3-.6-.6-.6zm-4.7 6.1c-.2 1.5-1.5 2.5-3 2.5-1.5 0-2.6-1.1-2.4-2.5.2-1.5 1.5-2.5 3-2.5 1.5 0 2.6 1.1 2.4 2.5z" fill="#003087"/>
              <path d="M37.8 9.8h-3.4c-.3 0-.6.2-.7.5l-3.9 8.6-1.7-8.5c-.1-.4-.4-.6-.7-.6h-3.4c-.4 0-.6.3-.5.7l3.6 13.5-3.5 4.9c-.3.4 0 1 .5 1h3.4c.3 0 .6-.2.8-.4l10.3-14.5c.3-.4 0-1-.5-1z" fill="#003087"/>
              <path d="M49.6 2.2H42c-.5 0-.9.4-1 .9L37.8 23.4c0 .3.2.6.5.6h3.8c.4 0 .8-.3.9-.7l.9-5.5c0-.3.3-.5.6-.5h2.2c4.4 0 7.7-1.8 8.6-6.6.4-2.1.1-3.8-.8-4.9-1.1-1.3-3-1.9-5.4-1.9zm1.8 6.2c-.4 2.1-1.8 3.4-3.8 3.4h-1.9c-.3 0-.6.2-.6.5l-.8 5.2c0 .3-.2.5-.5.5h-2.4c-.3 0-.5-.3-.5-.6l2.4-15.2c.1-.4.4-.7.8-.7h4.8c2.2 0 3.7.5 4.6 1.6.8 1.1.9 2.7.5 4.8z" fill="#0079C1"/>
              <path d="M60.4 9.8H57c-.3 0-.6.2-.6.5l-.2 1.1c-.5-.7-1.5-1.3-2.9-1.3-3.1 0-5.8 2.4-6.3 5.7-.3 1.7 0 3.3.9 4.5 1 1.2 2.5 1.8 4.2 1.8 2.5 0 3.9-1.6 3.9-1.6l-.2 1.1c0 .3.2.6.5.6h3.1c.4 0 .7-.3.7-.7l1.9-12.1c0-.4-.3-.6-.6-.6zm-4.7 6.1c-.2 1.5-1.5 2.5-3 2.5-1.5 0-2.6-1.1-2.4-2.5.2-1.5 1.5-2.5 3-2.5 1.5 0 2.6 1.1 2.4 2.5z" fill="#0079C1"/>
              <path d="M66.4 3.1l-2.6 16.7c0 .3.2.6.5.6h3.1c.4 0 .7-.3.7-.7l2.6-16.7c0-.3-.2-.6-.5-.6h-3.1c-.4 0-.7.3-.7.7z" fill="#0079C1"/>
            </svg>
          </div>

          <h3 className="text-base font-bold text-slate-900 mb-4">Pay with PayPal</h3>

          <div className="flex items-center gap-3.5 mb-5">
            <div className="w-10 h-10 rounded-xl bg-slate-50 border border-slate-200 flex items-center justify-center text-slate-600 flex-shrink-0">
              <Box className="w-5 h-5" />
            </div>
            <div>
              <div className="text-sm font-semibold text-slate-900">Apex Robotics Europe GmbH</div>
              <div className="text-xs text-slate-500">Industrial Servo Controller C101</div>
            </div>
          </div>

          <div className="flex justify-end items-baseline py-3 border-t border-slate-50 mb-5">
            <span className="text-xl font-bold text-slate-900 tracking-tight">$480.00 USD</span>
          </div>

          {/* Pay Button */}
          {demoState === 'idle' && (
            <button
              onClick={handleStartCheckout}
              className="w-full flex items-center justify-center gap-2 bg-[#0070BA] hover:bg-[#005EA6] active:scale-[0.98] text-white font-semibold py-3 px-5 rounded-full shadow-md shadow-[#0070BA]/20 transition-all cursor-pointer"
            >
              <Lock className="w-4 h-4" />
              <span>Pay Now</span>
            </button>
          )}

          {demoState === 'initiating' && (
            <button
              disabled
              className="w-full flex items-center justify-center gap-2 bg-[#0070BA]/80 text-white font-semibold py-3 px-5 rounded-full shadow-md cursor-not-allowed"
            >
              <Loader2 className="w-4 h-4 animate-spin" />
              <span>Contacting PayPal Sandbox...</span>
            </button>
          )}

          {demoState === 'timed_out' && (
            <button
              disabled
              className="w-full flex items-center justify-center gap-2 bg-red-600 text-white font-semibold py-3 px-5 rounded-full shadow-md shadow-red-600/20 cursor-not-allowed"
            >
              <span>HTTP 504 Gateway Timeout</span>
            </button>
          )}

          {(demoState === 'reconciling' || demoState === 'reconciled') && (
            <button
              disabled
              className="w-full flex items-center justify-center gap-2 bg-amber-600 text-white font-semibold py-3 px-5 rounded-full shadow-md shadow-amber-600/20 cursor-not-allowed"
            >
              <Lock className="w-4 h-4" />
              <span>DO_NOT_REPEAT_CHECKOUT (Protected)</span>
            </button>
          )}
        </div>

        {/* Right Card: PayBridge AI Recovery */}
        <div className="lg:col-span-7 bg-white border border-slate-200 rounded-2xl p-6 shadow-sm flex flex-col gap-4">
          
          {/* AI Header */}
          <div className="flex items-center justify-between pb-1">
            <div className="flex items-center gap-2.5">
              <div className="w-8 h-8 rounded-full bg-gradient-to-br from-violet-100 to-indigo-100 border border-indigo-200 flex items-center justify-center text-indigo-600">
                <Sparkles className="w-4 h-4" />
              </div>
              <span className="text-sm font-bold text-slate-900">PayBridge AI Recovery</span>
            </div>
            <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
              Live Simulation
            </span>
          </div>

          {/* Chat Bubble 1: Timeout Intercepted */}
          <div className="flex gap-2.5 items-start">
            <div className="w-6 h-6 rounded-md bg-slate-900 flex items-center justify-center text-white flex-shrink-0 mt-0.5">
              <Box className="w-3.5 h-3.5" />
            </div>
            <div className="flex-1 bg-slate-50 border border-slate-100 rounded-xl p-3 text-xs leading-relaxed text-slate-800">
              <span className="float-right text-[11px] text-slate-400 ml-2">Now</span>
              I detected a gateway timeout (HTTP 504) from PayPal. I'm checking the provider records to confirm your payment status.
            </div>
          </div>

          {/* Progress Box */}
          <div className="ml-8 bg-blue-50/60 border border-blue-100 rounded-xl p-3.5 space-y-2.5 text-xs font-medium text-slate-800">
            <div className="flex items-center gap-2">
              <CheckCircle2 className="w-4 h-4 text-emerald-600 flex-shrink-0" />
              <span>PayPal ledger verified</span>
            </div>
            <div className="flex items-center gap-2">
              <CheckCircle2 className="w-4 h-4 text-emerald-600 flex-shrink-0" />
              <span>Duplicate charge check (0 found)</span>
            </div>
            <div className="flex items-center gap-2 text-blue-600">
              {demoState === 'reconciled' ? (
                <>
                  <CheckCircle2 className="w-4 h-4 text-emerald-600 flex-shrink-0" />
                  <span className="text-slate-800 font-semibold">Reconciliation completed · Original capture verified</span>
                </>
              ) : (
                <>
                  <Loader2 className="w-4 h-4 animate-spin flex-shrink-0" />
                  <span>Reconciliation in progress...</span>
                </>
              )}
            </div>
          </div>

          {/* Shield Recommendation */}
          <div className="flex gap-2.5 items-start">
            <div className="w-6 h-6 rounded-full bg-violet-100 flex items-center justify-center text-violet-600 flex-shrink-0 mt-0.5">
              <ShieldAlert className="w-3.5 h-3.5" />
            </div>
            <div className="flex-1 bg-white border border-slate-200 rounded-xl p-3 shadow-xs">
              <div className="text-xs font-bold text-slate-900 mb-0.5">Do not repeat the payment.</div>
              <p className="text-[11.5px] text-slate-600">We're verifying your original payment and will update you shortly.</p>
            </div>
          </div>

          {/* User & AI Conversation Stream */}
          {chatHistory.length > 0 && (
            <div className="ml-8 space-y-2">
              {chatHistory.map((item, idx) => (
                <div key={idx} className="space-y-1.5">
                  <div className="ml-auto w-fit max-w-[85%] bg-[#0070BA] text-white rounded-xl px-3 py-1.5 text-xs">
                    {item.q}
                  </div>
                  <div className="mr-auto w-fit max-w-[95%] bg-slate-50 border border-slate-200 rounded-xl px-3 py-2 text-xs text-slate-800 leading-relaxed">
                    <strong>PayBridge AI:</strong> {item.a}
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* Chat Input & Suggested Chips */}
          <div className="mt-1 space-y-2">
            <div className="flex gap-1.5 flex-wrap">
              <button
                onClick={() => handleSendQuery('Was I charged?')}
                className="text-[11px] font-medium text-slate-600 bg-slate-100 hover:bg-slate-200 px-2.5 py-1 rounded-full transition-colors cursor-pointer"
              >
                Was I charged?
              </button>
              <button
                onClick={() => handleSendQuery("Why shouldn't I retry?")}
                className="text-[11px] font-medium text-slate-600 bg-slate-100 hover:bg-slate-200 px-2.5 py-1 rounded-full transition-colors cursor-pointer"
              >
                Why shouldn't I retry?
              </button>
              <button
                onClick={() => handleSendQuery('Where is my money?')}
                className="text-[11px] font-medium text-slate-600 bg-slate-100 hover:bg-slate-200 px-2.5 py-1 rounded-full transition-colors cursor-pointer"
              >
                Where is my money?
              </button>
            </div>

            <div className="flex items-center bg-white border border-slate-300 rounded-full px-3 py-1 focus-within:border-violet-500 focus-within:ring-2 focus-within:ring-violet-200 transition-all">
              <input
                type="text"
                value={inputText}
                onChange={e => setInputText(e.target.value)}
                onKeyDown={e => { if (e.key === 'Enter') handleSendQuery(); }}
                placeholder="Ask about this recovery..."
                className="flex-1 bg-transparent border-none outline-none text-xs text-slate-900 placeholder:text-slate-400 py-1"
              />
              <button
                onClick={() => handleSendQuery()}
                className="w-7 h-7 rounded-full bg-violet-600 hover:bg-violet-700 active:scale-95 text-white flex items-center justify-center transition-all cursor-pointer flex-shrink-0"
              >
                <Send className="w-3.5 h-3.5" />
              </button>
            </div>
          </div>

        </div>

      </div>

      {/* Section D: Final Result Outcome */}
      {demoState === 'reconciled' && (
        <div className="bg-emerald-50/70 border border-emerald-200 rounded-2xl p-5 animate-fade-in space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 text-sm font-bold text-emerald-800">
              <FileCheck2 className="w-5 h-5 text-emerald-600" />
              <span>Reconciliation Completed · Payment Confirmed ($0.00 Duplicate Exposure)</span>
            </div>
            <button
              onClick={handleReset}
              className="text-xs font-semibold px-3 py-1 rounded-md bg-white border border-slate-300 text-slate-700 hover:bg-slate-50 transition-colors cursor-pointer"
            >
              Reset Demo
            </button>
          </div>
          <p className="text-xs text-emerald-900">
            PayPal authoritatively confirms original capture of <strong>$480.00 USD</strong>. Customer was NOT billed twice. Cryptographic Truth Passport issued.
          </p>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 pt-1">
            <div className="bg-white border border-emerald-100 rounded-lg p-2.5">
              <div className="text-[10px] text-slate-500 uppercase tracking-wider">Provider Ref</div>
              <div className="text-xs font-mono font-bold text-slate-900 mt-0.5">Order #9K201824</div>
            </div>
            <div className="bg-white border border-emerald-100 rounded-lg p-2.5">
              <div className="text-[10px] text-slate-500 uppercase tracking-wider">Capture Status</div>
              <div className="text-xs font-mono font-bold text-slate-900 mt-0.5">cap-4819 (COMPLETED)</div>
            </div>
            <div className="bg-white border border-emerald-100 rounded-lg p-2.5">
              <div className="text-[10px] text-slate-500 uppercase tracking-wider">Duplicate Charges</div>
              <div className="text-xs font-mono font-bold text-slate-900 mt-0.5">$0.00 (Guaranteed)</div>
            </div>
            <div className="bg-white border border-emerald-100 rounded-lg p-2.5">
              <div className="text-[10px] text-slate-500 uppercase tracking-wider">Merkle Proof</div>
              <div className="text-xs font-mono font-bold text-slate-900 mt-0.5">4f8b91a2...7c9d</div>
            </div>
          </div>
        </div>
      )}

    </div>
  );
};
