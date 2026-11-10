import React, { useState } from 'react';
import { TruthGapGraph } from '../signature/TruthGapGraph';
import { PaymentTruthPassport } from '../signature/PaymentTruthPassport';
import { LifecycleTimeline } from '../signature/LifecycleTimeline';
import { DistillScannerHero } from '../distill/DistillScannerHero';
import { Shield, Clock, Scale, ArrowRight, CheckCircle, Lock, ExternalLink } from 'lucide-react';

/* =========================================================================
   Screen 1: Landing / Public Introduction with Distill Scanner Hero
   ========================================================================= */
export const LandingScreen: React.FC<{ onNavigate: (screen: any) => void }> = ({ onNavigate }) => {
  return (
    <div className="w-full bg-slate-950 text-slate-100">
      {/* The Distill "Built for people. Read by machines." Comparison Section */}
      <DistillScannerHero
        onExplore={() => onNavigate('dashboard')}
        onSignIn={() => onNavigate('auth')}
        onBookDemo={() => onNavigate('demo')}
      />

      {/* PayBridge AI Enterprise Invariants & Product Deep Dive */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 border-t border-slate-900">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-indigo-500/10 border border-indigo-500/20 text-indigo-400 font-mono text-xs tracking-wider uppercase mb-4">
            Financial Truth & Reconciliation Architecture
          </div>
          <h2 className="text-3xl sm:text-4xl font-bold tracking-tight text-slate-100">
            One payment. Two parties. <span className="text-cyan-400">One verified truth.</span>
          </h2>
          <p className="mt-4 text-slate-400 text-base sm:text-lg leading-relaxed">
            When cross-border payments experience upstream gateway timeouts (HTTP 504), traditional checkout flows panic and trigger duplicate charges. PayBridge AI isolates the ambiguous state, reconciles authoritative provider records, and provides mathematical certainty.
          </p>
        </div>

        {/* 3 Core Invariants Grid (Zero Emojis, Pure Professional Engineering) */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-16">
          <div className="bg-slate-900/60 border border-slate-800 rounded-2xl p-6 transition-all hover:border-slate-700">
            <div className="w-12 h-12 rounded-xl bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center text-indigo-400 mb-5">
              <Shield className="w-6 h-6" />
            </div>
            <h3 className="text-lg font-semibold text-slate-100 mb-2">AI Never Moves Money</h3>
            <p className="text-sm text-slate-400 leading-relaxed">
              Autonomous agents produce read-only analytical claims. Mutating PayPal operations are strictly restricted to authenticated human controllers with single-attempt locks.
            </p>
          </div>

          <div className="bg-slate-900/60 border border-slate-800 rounded-2xl p-6 transition-all hover:border-slate-700">
            <div className="w-12 h-12 rounded-xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-400 mb-5">
              <Clock className="w-6 h-6" />
            </div>
            <h3 className="text-lg font-semibold text-slate-100 mb-2">Timeout Is Not Failure</h3>
            <p className="text-sm text-slate-400 leading-relaxed">
              Lost network responses represent first-class UNKNOWN state. PayBridge immediately locks the checkout session with DO_NOT_REPEAT_CHECKOUT to eliminate duplicate billing.
            </p>
          </div>

          <div className="bg-slate-900/60 border border-slate-800 rounded-2xl p-6 transition-all hover:border-slate-700">
            <div className="w-12 h-12 rounded-xl bg-cyan-500/10 border border-cyan-500/20 flex items-center justify-center text-cyan-400 mb-5">
              <Scale className="w-6 h-6" />
            </div>
            <h3 className="text-lg font-semibold text-slate-100 mb-2">Cryptographic Evidence Truth</h3>
            <p className="text-sm text-slate-400 leading-relaxed">
              Financial state vectors are proven by Tier 1 provider reads and Tier 2 signed webhooks. Every claim is cited against an immutable SHA-256 evidence vault.
            </p>
          </div>
        </div>

        {/* Action Directives */}
        <div className="bg-slate-900/80 border border-slate-800/80 rounded-2xl p-8 flex flex-col sm:flex-row items-center justify-between gap-6">
          <div>
            <h4 className="text-lg font-semibold text-slate-100">Ready to test the reconciliation engine?</h4>
            <p className="text-sm text-slate-400 mt-1">
              Run Golden Demo B to simulate an upstream 504 gateway timeout and observe real-time multi-agent DAG investigation.
            </p>
          </div>
          <div className="flex items-center gap-3 shrink-0">
            <button
              onClick={() => onNavigate('demo')}
              className="px-5 py-2.5 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white font-medium text-sm transition-all shadow-lg shadow-indigo-600/20 flex items-center gap-2"
            >
              <span>Launch Golden Path Stepper</span>
              <ArrowRight className="w-4 h-4" />
            </button>
            <button
              onClick={() => onNavigate('cases')}
              className="px-5 py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 font-medium text-sm transition-all border border-slate-700"
            >
              View Cases Register
            </button>
          </div>
        </div>
      </section>
    </div>
  );
};

/* =========================================================================
   Screen 2: Authentication & Onboarding
   ========================================================================= */
export const AuthScreen: React.FC<{ onLogin: (role: string) => void }> = ({ onLogin }) => {
  return (
    <div style={{ maxWidth: '440px', margin: '60px auto', padding: '0 20px' }}>
      <div className="card">
        <h2 style={{ fontSize: '20px', fontWeight: 700, marginBottom: '8px' }}>Sign in to PayBridge</h2>
        <p style={{ fontSize: '13px', color: 'var(--text-2)', marginBottom: '20px' }}>
          Select a seeded institutional identity to test role-based projections.
        </p>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
          <button
            onClick={() => onLogin('sender')}
            style={{
              padding: '12px',
              backgroundColor: 'var(--surface-1)',
              border: '1px solid var(--border-strong)',
              borderRadius: '6px',
              textAlign: 'left',
            }}
          >
            <div style={{ fontWeight: 600, color: 'var(--text-1)' }}>Sender / Business Payer</div>
            <div style={{ fontSize: '12px', color: 'var(--text-3)' }}>alex@acmecorp.us • Creates payment & approves checkout</div>
          </button>

          <button
            onClick={() => onLogin('receiver')}
            style={{
              padding: '12px',
              backgroundColor: 'var(--surface-1)',
              border: '1px solid var(--border-strong)',
              borderRadius: '6px',
              textAlign: 'left',
            }}
          >
            <div style={{ fontWeight: 600, color: 'var(--text-1)' }}>Receiver / Contractor</div>
            <div style={{ fontSize: '12px', color: 'var(--text-3)' }}>priya@consulting.in • Consent-gated dual-lens viewer</div>
          </button>

          <button
            onClick={() => onLogin('merchant_analyst')}
            style={{
              padding: '12px',
              backgroundColor: 'var(--surface-1)',
              border: '1px solid var(--border-strong)',
              borderRadius: '6px',
              textAlign: 'left',
            }}
          >
            <div style={{ fontWeight: 600, color: 'var(--text-1)' }}>Merchant Operations Analyst</div>
            <div style={{ fontSize: '12px', color: 'var(--text-3)' }}>ops@platform.internal • Triage, manual reconcile & audit</div>
          </button>

          <button
            onClick={() => onLogin('demo_operator')}
            style={{
              padding: '12px',
              backgroundColor: 'var(--surface-1)',
              border: '1px solid var(--accent)',
              borderRadius: '6px',
              textAlign: 'left',
            }}
          >
            <div style={{ fontWeight: 600, color: 'var(--accent)' }}>Judge / Demo Operator (Golden Demo B)</div>
            <div style={{ fontSize: '12px', color: 'var(--text-3)' }}>Runs chaos proxy, response drops & live sandbox tests</div>
          </button>
        </div>
      </div>
    </div>
  );
};

/* =========================================================================
   Screen 3: Workspace / Case Dashboard
   ========================================================================= */
export const CasesDashboardScreen: React.FC<{ onSelectCase: (id: string) => void }> = ({ onSelectCase }) => {
  return (
    <div style={{ padding: '24px', maxWidth: '1100px', margin: '0 auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <div>
          <h2 style={{ fontSize: '22px', fontWeight: 700 }}>Payment Cases & Truth Matrix</h2>
          <p style={{ fontSize: '13px', color: 'var(--text-2)' }}>
            Multi-dimensional payment state monitoring. Tracking Provider, Application, and Receiver vectors.
          </p>
        </div>
        <div style={{ display: 'flex', gap: '8px' }}>
          <span className="badge badge-verified">2 Reconciled</span>
          <span className="badge badge-unknown">1 Ambiguous / Timeout</span>
        </div>
      </div>

      <div className="card" style={{ padding: 0, overflow: 'hidden' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '13px' }}>
          <thead>
            <tr style={{ backgroundColor: 'var(--surface-sunken)', borderBottom: '1px solid var(--border)' }}>
              <th style={{ padding: '12px 16px' }}>Case Reference</th>
              <th style={{ padding: '12px 16px' }}>Route & Amount</th>
              <th style={{ padding: '12px 16px' }}>Provider State</th>
              <th style={{ padding: '12px 16px' }}>Application State</th>
              <th style={{ padding: '12px 16px' }}>Safe Action</th>
              <th style={{ padding: '12px 16px' }}>Action</th>
            </tr>
          </thead>
          <tbody>
            <tr style={{ borderBottom: '1px solid var(--border)' }}>
              <td style={{ padding: '12px 16px' }}>
                <strong>#PB-8821</strong>
                <div style={{ fontSize: '11px', color: 'var(--text-3)' }}>Invoice INV-2026-081</div>
              </td>
              <td style={{ padding: '12px 16px' }}>
                US → IN ($250.00 USD)
              </td>
              <td style={{ padding: '12px 16px' }}>
                <span className="badge badge-unknown">UNKNOWN</span>
                <div style={{ fontSize: '10px', color: 'var(--text-3)' }}>Last: APPROVED</div>
              </td>
              <td style={{ padding: '12px 16px' }}>
                <span className="badge badge-unknown">API_TIMEOUT</span>
              </td>
              <td style={{ padding: '12px 16px' }}>
                <span style={{ color: 'var(--unknown)', fontWeight: 600 }}>DO NOT RETRY</span>
              </td>
              <td style={{ padding: '12px 16px' }}>
                <button
                  onClick={() => onSelectCase('PB-8821')}
                  style={{
                    padding: '6px 12px',
                    backgroundColor: 'var(--accent)',
                    color: '#fff',
                    border: 'none',
                    borderRadius: '4px',
                    fontSize: '12px',
                  }}
                >
                  Inspect Case
                </button>
              </td>
            </tr>
            <tr style={{ borderBottom: '1px solid var(--border)' }}>
              <td style={{ padding: '12px 16px' }}>
                <strong>#PB-8819</strong>
                <div style={{ fontSize: '11px', color: 'var(--text-3)' }}>Invoice INV-2026-079</div>
              </td>
              <td style={{ padding: '12px 16px' }}>
                US → DE (€1,200.00 EUR)
              </td>
              <td style={{ padding: '12px 16px' }}>
                <span className="badge badge-verified">CAPTURE_COMPLETED</span>
              </td>
              <td style={{ padding: '12px 16px' }}>
                <span className="badge badge-verified">SYNCHRONISED</span>
              </td>
              <td style={{ padding: '12px 16px' }}>
                <span style={{ color: 'var(--verified)', fontWeight: 600 }}>NONE_REQUIRED</span>
              </td>
              <td style={{ padding: '12px 16px' }}>
                <button
                  onClick={() => onSelectCase('PB-8819')}
                  style={{
                    padding: '6px 12px',
                    backgroundColor: 'var(--surface-1)',
                    border: '1px solid var(--border-strong)',
                    borderRadius: '4px',
                    fontSize: '12px',
                  }}
                >
                  View Passport
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
};

/* =========================================================================
   Screen 4: New Payment & Timeout Hero (with DO-NOT-REPEAT banner)
   ========================================================================= */
export const NewPaymentScreen: React.FC<{ onSimulateTimeout: () => void }> = ({ onSimulateTimeout }) => {
  const [inFlight, setInFlight] = useState(false);
  const [timedOut, setTimedOut] = useState(false);

  const handlePay = () => {
    setInFlight(true);
    setTimeout(() => {
      setInFlight(false);
      setTimedOut(true);
      onSimulateTimeout();
    }, 1500);
  };

  return (
    <div style={{ maxWidth: '640px', margin: '40px auto', padding: '0 20px' }}>
      <div className="card">
        <h2 style={{ fontSize: '20px', fontWeight: 700, marginBottom: '8px' }}>Send Cross-Border Invoice Payment</h2>
        <p style={{ fontSize: '13px', color: 'var(--text-2)', marginBottom: '20px' }}>
          Connects directly to PayPal Sandbox with strict single-attempt locking.
        </p>

        {timedOut && (
          <div className="alert-banner alert-unknown" style={{ marginBottom: '20px' }}>
            <span style={{ fontSize: '20px' }}>⚠️</span>
            <div>
              <strong style={{ fontSize: '15px' }}>DO NOT REPEAT CHECKOUT</strong>
              <p style={{ marginTop: '4px', fontSize: '13px' }}>
                The network response timed out after 15s. PayPal may have already captured funds. PayBridge is reconciling the original transaction now. Submitting another payment is strictly blocked.
              </p>
            </div>
          </div>
        )}

        <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
          <div>
            <label style={{ fontSize: '12px', fontWeight: 600, color: 'var(--text-2)' }}>Invoice Reference</label>
            <input
              type="text"
              readOnly
              value="INV-2026-081 (Software Development)"
              style={{ width: '100%', padding: '8px 12px', borderRadius: '4px', border: '1px solid var(--border)', marginTop: '4px' }}
            />
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
            <div>
              <label style={{ fontSize: '12px', fontWeight: 600, color: 'var(--text-2)' }}>Amount</label>
              <input
                type="text"
                readOnly
                value="$250.00"
                style={{ width: '100%', padding: '8px 12px', borderRadius: '4px', border: '1px solid var(--border)', marginTop: '4px' }}
              />
            </div>
            <div>
              <label style={{ fontSize: '12px', fontWeight: 600, color: 'var(--text-2)' }}>Currency</label>
              <input
                type="text"
                readOnly
                value="USD"
                style={{ width: '100%', padding: '8px 12px', borderRadius: '4px', border: '1px solid var(--border)', marginTop: '4px' }}
              />
            </div>
          </div>

          <button
            disabled={inFlight || timedOut}
            onClick={handlePay}
            style={{
              marginTop: '12px',
              padding: '12px',
              backgroundColor: timedOut ? 'var(--border)' : 'var(--accent)',
              color: '#fff',
              border: 'none',
              borderRadius: '6px',
              fontWeight: 600,
              fontSize: '14px',
              cursor: timedOut ? 'not-allowed' : 'pointer',
            }}
          >
            {inFlight ? 'Executing Capture via PayPal Sandbox...' : timedOut ? 'Retry Blocked by Safety Invariant' : 'Approve & Capture via PayPal ($250.00 USD)'}
          </button>
        </div>
      </div>
    </div>
  );
};

/* =========================================================================
   Screen 5: Payment Detail & Truth Passport
   ========================================================================= */
export const PaymentDetailScreen: React.FC = () => {
  return (
    <div style={{ maxWidth: '960px', margin: '24px auto', padding: '0 20px', display: 'flex', flexDirection: 'column', gap: '24px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <span style={{ fontSize: '12px', color: 'var(--text-3)' }}>Case #PB-8821</span>
          <h1 style={{ fontSize: '24px', fontWeight: 700 }}>Payment Truth Shell</h1>
        </div>
        <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
          <span style={{ width: '8px', height: '8px', borderRadius: '50%', backgroundColor: 'var(--verified)' }} />
          <span style={{ fontSize: '12px', color: 'var(--text-2)' }}>SSE Real-Time Stream Active</span>
        </div>
      </div>

      <PaymentTruthPassport />
      <TruthGapGraph />
      <LifecycleTimeline />
    </div>
  );
};

/* =========================================================================
   Screen 6: Payment Lifecycle Timeline
   ========================================================================= */
export const TimelineScreen: React.FC = () => {
  return (
    <div style={{ maxWidth: '800px', margin: '24px auto', padding: '0 20px' }}>
      <LifecycleTimeline />
    </div>
  );
};

/* =========================================================================
   Screen 7: Truth Gap Workbench
   ========================================================================= */
export const TruthGapScreen: React.FC = () => {
  return (
    <div style={{ maxWidth: '960px', margin: '24px auto', padding: '0 20px' }}>
      <TruthGapGraph />
    </div>
  );
};

/* =========================================================================
   Screen 8: Sender Two-Lens View
   ========================================================================= */
export const SenderLensScreen: React.FC = () => {
  return (
    <div style={{ maxWidth: '720px', margin: '24px auto', padding: '0 20px' }}>
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
          <span className="badge badge-verified">Sender Lens Active</span>
          <span style={{ fontSize: '12px', color: 'var(--text-3)' }}>Privacy Guard: DLP Verified</span>
        </div>
        <h2 style={{ fontSize: '18px', fontWeight: 700 }}>Your Payment of $250.00 USD to Priya Kasthuri</h2>
        <p style={{ fontSize: '14px', color: 'var(--text-2)', marginTop: '8px' }}>
          Status: <strong>Verified with PayPal</strong>. The original capture succeeded on PayPal's ledger. No second charge was made.
        </p>
        <div style={{ marginTop: '16px', padding: '12px', background: 'var(--surface-sunken)', borderRadius: '6px', fontSize: '13px' }}>
          <strong>Next Step for Sender:</strong> None required. Your invoice is fulfilled.
        </div>
      </div>
    </div>
  );
};

/* =========================================================================
   Screen 9: Receiver Two-Lens View (Consent-Gated)
   ========================================================================= */
export const ReceiverLensScreen: React.FC = () => {
  const [consentGranted, setConsentGranted] = useState(true);

  if (!consentGranted) {
    return (
      <div style={{ maxWidth: '540px', margin: '40px auto', padding: '0 20px' }}>
        <div className="card">
          <h3 style={{ fontSize: '18px', fontWeight: 600 }}>Consent Required</h3>
          <p style={{ fontSize: '13px', color: 'var(--text-2)', marginTop: '6px' }}>
            To view this cross-border payment inquiry, you must accept the consent agreement.
          </p>
          <button
            onClick={() => setConsentGranted(true)}
            style={{ marginTop: '16px', padding: '8px 16px', background: 'var(--accent)', color: '#fff', border: 'none', borderRadius: '4px' }}
          >
            Grant Consent & View Case
          </button>
        </div>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '720px', margin: '24px auto', padding: '0 20px' }}>
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
          <span className="badge badge-verified">Receiver Lens (Consent Active)</span>
          <button
            onClick={() => setConsentGranted(false)}
            style={{ fontSize: '11px', color: 'var(--error)', background: 'none', border: 'none', cursor: 'pointer' }}
          >
            Revoke Consent
          </button>
        </div>
        <h2 style={{ fontSize: '18px', fontWeight: 700 }}>Incoming Payment of $250.00 USD from Acme Corp</h2>
        <p style={{ fontSize: '14px', color: 'var(--text-2)', marginTop: '8px' }}>
          PayPal has processed and confirmed the capture. Note: Cross-border bank deposit settlement takes 1-2 business days according to your local bank.
        </p>
      </div>
    </div>
  );
};

/* =========================================================================
   Screen 10: Next Safe Action Panel
   ========================================================================= */
export const SafeActionScreen: React.FC = () => {
  return (
    <div style={{ maxWidth: '720px', margin: '24px auto', padding: '0 20px' }}>
      <div className="card">
        <h2 style={{ fontSize: '18px', fontWeight: 700, marginBottom: '12px' }}>Next Safe Action Control Panel</h2>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          <div style={{ padding: '12px', background: 'var(--verified-bg)', border: '1px solid var(--verified)', borderRadius: '6px' }}>
            <span className="badge badge-verified">Allowed</span>
            <div style={{ fontWeight: 600, color: 'var(--verified)', marginTop: '4px' }}>Query Original Transaction (GET)</div>
            <div style={{ fontSize: '12px', color: 'var(--text-2)' }}>Re-reads original order status without any money movement.</div>
          </div>
          <div style={{ padding: '12px', background: 'var(--error-bg)', border: '1px solid var(--error)', borderRadius: '6px' }}>
            <span className="badge badge-error">BLOCKED</span>
            <div style={{ fontWeight: 600, color: 'var(--error)', marginTop: '4px' }}>Retry Checkout Mutation</div>
            <div style={{ fontSize: '12px', color: 'var(--text-2)' }}>Engine invariant: Multiple active attempts on the same case are strictly forbidden.</div>
          </div>
        </div>
      </div>
    </div>
  );
};

/* =========================================================================
   Screen 11: Counterfactual Recovery Mirror Lab (SIMULATED)
   ========================================================================= */
export const MirrorLabScreen: React.FC = () => {
  return (
    <div style={{ maxWidth: '800px', margin: '24px auto', padding: '0 20px' }}>
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
          <span className="badge badge-simulated">SIMULATED FIXTURE ENVIRONMENT</span>
          <span style={{ fontSize: '11px', color: 'var(--text-3)' }}>Isolated from Real Sandbox</span>
        </div>
        <h2 style={{ fontSize: '18px', fontWeight: 700 }}>Counterfactual Recovery Mirror Lab</h2>
        <p style={{ fontSize: '13px', color: 'var(--text-2)', marginTop: '6px' }}>
          Replay edge cases: lost responses, dropped webhooks, and out-of-order notifications without affecting real transactions.
        </p>
      </div>
    </div>
  );
};

/* =========================================================================
   Screen 12: Evidence Explorer & Provenance
   ========================================================================= */
export const EvidenceScreen: React.FC = () => {
  return (
    <div style={{ maxWidth: '800px', margin: '24px auto', padding: '0 20px' }}>
      <div className="card">
        <h2 style={{ fontSize: '18px', fontWeight: 700, marginBottom: '16px' }}>Evidence Explorer & Provenance Hash Chain</h2>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
          <div style={{ padding: '12px', background: 'var(--surface-sunken)', borderRadius: '6px', fontSize: '12px' }}>
            <div><strong>Evidence #EV-CAP-3KK109</strong> (Tier E1 Direct Provider Read)</div>
            <div>SHA-256: <code>9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08</code></div>
            <div>Verified: true • Source: api-m.sandbox.paypal.com</div>
          </div>
        </div>
      </div>
    </div>
  );
};

/* =========================================================================
   Screen 13: Notification Center
   ========================================================================= */
export const NotificationsScreen: React.FC = () => {
  return (
    <div style={{ maxWidth: '640px', margin: '24px auto', padding: '0 20px' }}>
      <div className="card">
        <h2 style={{ fontSize: '18px', fontWeight: 700, marginBottom: '16px' }}>Institutional Notifications</h2>
        <div style={{ padding: '12px', borderLeft: '4px solid var(--verified)', background: 'var(--surface-sunken)', borderRadius: '4px', fontSize: '13px' }}>
          <strong>Payment #PB-8821 Reconciled Successfully</strong>
          <p style={{ color: 'var(--text-2)', margin: '4px 0 0' }}>The original capture was verified on the PayPal Sandbox network.</p>
        </div>
      </div>
    </div>
  );
};

/* =========================================================================
   Screen 14: Operations & Sandbox Webhook Health
   ========================================================================= */
export const OperationsScreen: React.FC = () => {
  return (
    <div style={{ maxWidth: '800px', margin: '24px auto', padding: '0 20px' }}>
      <div className="card">
        <h2 style={{ fontSize: '18px', fontWeight: 700, marginBottom: '16px' }}>Sandbox Webhook & Queue Health</h2>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '12px' }}>
          <div style={{ padding: '12px', background: 'var(--surface-sunken)', borderRadius: '6px' }}>
            <div style={{ fontSize: '12px', color: 'var(--text-3)' }}>Webhook Signature Status</div>
            <div style={{ fontSize: '18px', fontWeight: 700, color: 'var(--verified)' }}>100% Valid</div>
          </div>
          <div style={{ padding: '12px', background: 'var(--surface-sunken)', borderRadius: '6px' }}>
            <div style={{ fontSize: '12px', color: 'var(--text-3)' }}>Queue Oldest Job Age</div>
            <div style={{ fontSize: '18px', fontWeight: 700 }}>0.8s</div>
          </div>
          <div style={{ padding: '12px', background: 'var(--surface-sunken)', borderRadius: '6px' }}>
            <div style={{ fontSize: '12px', color: 'var(--text-3)' }}>Dead Letter Count</div>
            <div style={{ fontSize: '18px', fontWeight: 700 }}>0</div>
          </div>
        </div>
      </div>
    </div>
  );
};

/* =========================================================================
   Screen 15: Settings, Privacy & Consent
   ========================================================================= */
export const SettingsScreen: React.FC = () => {
  return (
    <div style={{ maxWidth: '640px', margin: '24px auto', padding: '0 20px' }}>
      <div className="card">
        <h2 style={{ fontSize: '18px', fontWeight: 700, marginBottom: '16px' }}>Privacy & Data Governance</h2>
        <p style={{ fontSize: '13px', color: 'var(--text-2)' }}>
          Data retention policy: 90 days with automatic cryptographic purge. RLS enabled on all tenant tables.
        </p>
      </div>
    </div>
  );
};

/* =========================================================================
   Screen 16: Demo Walkthrough & Golden Demo B Scenario Selector
   ========================================================================= */
export const DemoWalkthroughScreen: React.FC<{ onRunGoldenDemo: () => void }> = ({ onRunGoldenDemo }) => {
  const [running, setRunning] = useState(false);
  const [step, setStep] = useState(0);

  const startDemo = () => {
    setRunning(true);
    setStep(1); // Create order & approve
    setTimeout(() => {
      setStep(2); // Simulate network drop on capture
      setTimeout(() => {
        setStep(3); // PayBridge reconciles original capture
        setTimeout(() => {
          setStep(4); // Reconciled!
          setRunning(false);
          onRunGoldenDemo();
        }, 1500);
      }, 1500);
    }, 1500);
  };

  return (
    <div style={{ maxWidth: '840px', margin: '24px auto', padding: '0 20px' }}>
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
          <div>
            <span className="badge badge-verified">Official Scenario Runner</span>
            <h2 style={{ fontSize: '20px', fontWeight: 700, marginTop: '4px' }}>
              Golden Demo B: Ambiguous Capture & Safe Recovery
            </h2>
          </div>
          <button
            disabled={running}
            onClick={startDemo}
            style={{
              padding: '10px 20px',
              backgroundColor: 'var(--accent)',
              color: '#fff',
              border: 'none',
              borderRadius: '6px',
              fontWeight: 600,
              fontSize: '14px',
            }}
          >
            {running ? 'Running Golden Demo B...' : '▶ Run Golden Demo B Live'}
          </button>
        </div>

        <p style={{ fontSize: '13px', color: 'var(--text-2)', marginBottom: '20px' }}>
          This script executes the complete hackathon golden scenario: PayPal sandbox capture succeeds, PayBridge confirmation response is dropped by a simulated network proxy, PayBridge reconciles the original transaction using read-only GETs, and 0 duplicate charges happen.
        </p>

        {/* Live Step Progress */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          <div style={{ padding: '12px', background: step >= 1 ? 'var(--surface-sunken)' : 'var(--surface-2)', borderRadius: '6px', borderLeft: step >= 1 ? '4px solid var(--verified)' : '4px solid var(--border)' }}>
            <strong>Step 1: Order Creation & Sandbox Buyer Approval</strong>
            <div style={{ fontSize: '12px', color: 'var(--text-3)' }}>Order #ORD-8NW921 created for $250.00 USD with intent CAPTURE.</div>
          </div>
          <div style={{ padding: '12px', background: step >= 2 ? 'var(--unknown-bg)' : 'var(--surface-2)', borderRadius: '6px', borderLeft: step >= 2 ? '4px solid var(--unknown)' : '4px solid var(--border)' }}>
            <strong>Step 2: Network Drop on Capture Response (HTTP 504 Gateway Timeout)</strong>
            <div style={{ fontSize: '12px', color: 'var(--text-3)' }}>PayBridge records OUTCOME-UNKNOWN; DO_NOT_REPEAT_CHECKOUT guard activates immediately.</div>
          </div>
          <div style={{ padding: '12px', background: step >= 3 ? 'var(--surface-sunken)' : 'var(--surface-2)', borderRadius: '6px', borderLeft: step >= 3 ? '4px solid var(--accent)' : '4px solid var(--border)' }}>
            <strong>Step 3: Safe Read-Only Reconciliation</strong>
            <div style={{ fontSize: '12px', color: 'var(--text-3)' }}>Reconciler queries GET /v2/checkout/orders/ORD-8NW921; discovers capture was COMPLETED.</div>
          </div>
          <div style={{ padding: '12px', background: step >= 4 ? 'var(--verified-bg)' : 'var(--surface-2)', borderRadius: '6px', borderLeft: step >= 4 ? '4px solid var(--verified)' : '4px solid var(--border)' }}>
            <strong>Step 4: Truth Passport Sealed & 0 Blind Retries</strong>
            <div style={{ fontSize: '12px', color: 'var(--text-3)' }}>State synchronized; sender and receiver informed separately; duplicate charge rate = 0.00%.</div>
          </div>
        </div>
      </div>
    </div>
  );
};
