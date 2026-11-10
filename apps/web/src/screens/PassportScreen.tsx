import React, { useState } from 'react';
import { 
  Award, 
  ShieldCheck, 
  Lock, 
  FileCheck, 
  Eye, 
  Scale, 
  CheckCircle2, 
  Copy, 
  ExternalLink,
  QrCode
} from 'lucide-react';
import { StatusPill } from '../components/StatusPill';

export const PassportScreen: React.FC = () => {
  const [lens, setLens] = useState<'merchant' | 'investigator'>('merchant');
  const [copied, setCopied] = useState(false);

  const passportHash = '4fc20f0e89261a829b68f023fcc582934dae262520ded004aa12ec660191ac6b';

  const copyHash = () => {
    navigator.clipboard.writeText(passportHash);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="space-y-8 animate-fade-in max-w-5xl mx-auto pb-12">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <Award className="w-6 h-6 text-indigo-400" />
            <h1 className="text-2xl font-bold tracking-tight text-slate-100">
              Payment Truth Passport
            </h1>
          </div>
          <p className="text-sm text-slate-400 mt-1">
            Authoritative, cryptographically verifiable financial settlement passport.
          </p>
        </div>

        {/* Lens Switcher */}
        <div className="flex items-center p-1 bg-slate-900 border border-slate-800 rounded-xl">
          <button
            onClick={() => setLens('merchant')}
            className={`flex items-center gap-2 px-3.5 py-1.5 rounded-lg text-xs font-medium transition-all cursor-pointer ${
              lens === 'merchant'
                ? 'bg-indigo-600 text-white shadow-sm'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            <Eye className="w-3.5 h-3.5" />
            Merchant Lens
          </button>
          <button
            onClick={() => setLens('investigator')}
            className={`flex items-center gap-2 px-3.5 py-1.5 rounded-lg text-xs font-medium transition-all cursor-pointer ${
              lens === 'investigator'
                ? 'bg-cyan-600 text-white shadow-sm'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            <Scale className="w-3.5 h-3.5" />
            Investigator Lens
          </button>
        </div>
      </div>

      {/* Passport Certificate Card */}
      <div className="rounded-2xl bg-gradient-to-b from-slate-900 via-slate-900 to-slate-950 border border-indigo-500/30 shadow-2xl p-6 sm:p-8 relative overflow-hidden">
        {/* Background Watermark */}
        <div className="absolute right-0 bottom-0 opacity-5 pointer-events-none transform translate-x-12 translate-y-12">
          <ShieldCheck className="w-96 h-96 text-indigo-400" />
        </div>

        {/* Certificate Top Ribbon */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-6 border-b border-slate-800 gap-4">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 rounded-xl bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center">
              <ShieldCheck className="w-7 h-7 text-indigo-400" />
            </div>
            <div>
              <span className="text-[10px] font-mono tracking-widest uppercase font-bold text-indigo-400">
                OFFICIAL FINANCIAL STATE PROOF
              </span>
              <h2 className="text-xl font-bold text-slate-100">
                Payment Truth Passport v2
              </h2>
            </div>
          </div>

          <div className="text-right">
            <span className="text-xs text-slate-400 font-mono">Passport State Hash:</span>
            <div className="flex items-center gap-2 mt-0.5">
              <code className="text-xs font-mono text-cyan-300 bg-slate-950 px-2.5 py-1 rounded border border-slate-800">
                {passportHash.slice(0, 20)}...
              </code>
              <button 
                onClick={copyHash}
                title="Copy full cryptographic SHA-256 hash"
                className="p-1 rounded hover:bg-slate-800 text-slate-400 hover:text-slate-200 cursor-pointer"
              >
                <Copy className="w-3.5 h-3.5" />
              </button>
            </div>
            {copied && <span className="text-[10px] text-emerald-400">Copied to clipboard!</span>}
          </div>
        </div>

        {/* Lens Explanation Banner */}
        <div className="my-6 p-4 rounded-xl bg-slate-950/70 border border-slate-800/80">
          <div className="flex items-center gap-2 mb-1.5">
            <span className="text-xs font-mono font-bold uppercase text-indigo-400">
              {lens === 'merchant' ? 'Operational Merchant Perspective' : 'Cryptographic Investigator Perspective'}
            </span>
          </div>

          {lens === 'merchant' ? (
            <p className="text-sm text-slate-200 leading-relaxed">
              <strong>PayPal confirmed your payment capture.</strong> Your website did not receive the initial confirmation due to an upstream gateway timeout. PayBridge automatically matched the capture with invoice <strong>INV-2026-8941</strong> and reconciled the financial state. <em>Please do not charge the customer again.</em>
            </p>
          ) : (
            <div className="text-xs font-mono text-slate-300 space-y-1">
              <div>• Audit Chain Head: <span className="text-cyan-400">0x7a89f...21c</span> (Transactionally Locked)</div>
              <div>• Signed Assertion: HMAC-SHA256 Token with 60s TTL verified against PostgreSQL GUCs</div>
              <div>• Canonical Argument Hash (RFC 8785): <span className="text-emerald-400">4fc20f0e89261a829b68f023fcc582934dae262520ded004aa12ec660191ac6b</span></div>
              <div>• Evidence Ingestion: Tier 1 PROVIDER_VERIFIED payload stored immutably with SHA-256 digest</div>
            </div>
          )}
        </div>

        {/* 3-Pillar Triad Comparison */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-5 pt-2">
          {/* Pillar 1: Sender Claim */}
          <div className="p-4 rounded-xl bg-slate-950/50 border border-slate-800 space-y-3">
            <div className="flex items-center justify-between border-b border-slate-800/80 pb-2">
              <span className="text-xs font-bold text-slate-300 font-mono">1. SENDER INTENT</span>
              <StatusPill status="SYNCHRONISED" />
            </div>
            <div className="text-xs space-y-1.5 text-slate-300 font-mono">
              <div>Invoice: <strong className="text-slate-100">INV-2026-8941</strong></div>
              <div>Amount: <strong className="text-slate-100">$149.50 USD</strong></div>
              <div>Purpose: <span>CONSULTING_SERVICES</span></div>
              <div>Corridor: <span>US → DE</span></div>
              <div>Buyer Status: <span className="text-emerald-400">Paid & Approved</span></div>
            </div>
          </div>

          {/* Pillar 2: Authoritative Provider Truth */}
          <div className="p-4 rounded-xl bg-slate-950/50 border border-indigo-500/30 space-y-3 shadow-inner">
            <div className="flex items-center justify-between border-b border-slate-800/80 pb-2">
              <span className="text-xs font-bold text-indigo-300 font-mono">2. PROVIDER GROUND TRUTH</span>
              <StatusPill status="CAPTURE_COMPLETED" />
            </div>
            <div className="text-xs space-y-1.5 text-slate-300 font-mono">
              <div>Capture ID: <strong className="text-indigo-300">CAP-PAYPAL-782199</strong></div>
              <div>Captured: <strong className="text-emerald-400">$149.50 USD</strong></div>
              <div>Seller Protect: <span className="text-emerald-400">ELIGIBLE</span></div>
              <div>Authoritative GET: <span className="text-cyan-400">200 OK</span></div>
              <div>Dispute Risk: <span className="text-emerald-400">0% (Truth Established)</span></div>
            </div>
          </div>

          {/* Pillar 3: Receiver Settlement */}
          <div className="p-4 rounded-xl bg-slate-950/50 border border-slate-800 space-y-3">
            <div className="flex items-center justify-between border-b border-slate-800/80 pb-2">
              <span className="text-xs font-bold text-slate-300 font-mono">3. RECEIVER SETTLEMENT</span>
              <StatusPill status="CREDITED" />
            </div>
            <div className="text-xs space-y-1.5 text-slate-300 font-mono">
              <div>Receiver Account: <span className="text-slate-200">merchant@seller.de</span></div>
              <div>Gross Received: <strong className="text-slate-100">$149.50 USD</strong></div>
              <div>State: <span className="text-emerald-400">AVAILABLE</span></div>
              <div>Bank Deposit: <span className="text-amber-400">NOT_OBSERVABLE</span></div>
              <div>Reconciled: <span className="text-emerald-400">YES</span></div>
            </div>
          </div>
        </div>

        {/* Certificate Bottom Verification Stamp */}
        <div className="mt-8 pt-6 border-t border-slate-800 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-lg bg-emerald-500/10 border border-emerald-500/30 flex items-center justify-center">
              <CheckCircle2 className="w-5 h-5 text-emerald-400" />
            </div>
            <div className="text-xs text-slate-400">
              <span className="font-semibold text-slate-200">Cryptographically Sealed by PayBridge AI Engine</span>
              <div>Deterministic state verified under PostgreSQL 16 row-level security policies.</div>
            </div>
          </div>

          <button
            onClick={() => window.print()}
            className="flex items-center gap-2 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium rounded-lg border border-slate-700 transition-all cursor-pointer"
          >
            <ExternalLink className="w-3.5 h-3.5" />
            Export Official PDF Proof
          </button>
        </div>
      </div>
    </div>
  );
};
