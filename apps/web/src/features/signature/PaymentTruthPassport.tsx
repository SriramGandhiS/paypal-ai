import React from 'react';
import { PaymentTruthPassport as PassportType } from '../../types';

interface PaymentTruthPassportProps {
  passport?: PassportType;
}

export const PaymentTruthPassport: React.FC<PaymentTruthPassportProps> = ({
  passport = defaultPassport,
}) => {
  return (
    <div className="card" style={{ border: '2px solid var(--border-strong)', background: 'var(--surface-1)' }}>
      {/* Passport Header */}
      <div
        style={{
          borderBottom: '1px solid var(--border)',
          paddingBottom: '16px',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'flex-start',
        }}
      >
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '4px' }}>
            <span style={{ fontSize: '20px' }}>📜</span>
            <h2 style={{ fontSize: '20px', fontWeight: 700, color: 'var(--text-1)' }}>
              Payment Truth Passport
            </h2>
            <span className="badge badge-verified">v{passport.version} Sealed</span>
          </div>
          <p style={{ fontSize: '13px', color: 'var(--text-2)' }}>
            Authoritative, evidence-grounded financial truth document. Cites 100% verified provider records.
          </p>
        </div>
        <div style={{ textAlign: 'right', fontSize: '11px', color: 'var(--text-3)' }}>
          <div>SHA-256 Fingerprint:</div>
          <code style={{ fontSize: '10px', color: 'var(--text-1)' }}>
            {passport.content_sha256.substring(0, 16)}...
          </code>
        </div>
      </div>

      {/* Confirmed Facts Section */}
      <div style={{ marginTop: '20px' }}>
        <h4 style={{ fontSize: '14px', textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--text-2)', marginBottom: '12px' }}>
          Confirmed Financial Facts (Provider Evidenced)
        </h4>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
          {passport.confirmed_facts.map((fact, idx) => (
            <div
              key={idx}
              style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                padding: '10px 14px',
                background: 'var(--surface-2)',
                borderRadius: '6px',
                borderLeft: '4px solid var(--verified)',
              }}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                <span className="badge badge-verified">Tier E{fact.tier}</span>
                <span style={{ fontSize: '14px', fontWeight: 500, color: 'var(--text-1)' }}>
                  {fact.statement}
                </span>
              </div>
              <code style={{ fontSize: '11px', color: 'var(--text-3)' }}>
                Evidence: {fact.evidence_id}
              </code>
            </div>
          ))}
        </div>
      </div>

      {/* Uncertain Facts Section */}
      {passport.uncertain_facts.length > 0 && (
        <div style={{ marginTop: '20px' }}>
          <h4 style={{ fontSize: '14px', textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--unknown)', marginBottom: '12px' }}>
            Uncertain or Unobserved Boundary Facts
          </h4>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
            {passport.uncertain_facts.map((u, idx) => (
              <div
                key={idx}
                style={{
                  padding: '10px 14px',
                  background: 'var(--unknown-bg)',
                  borderRadius: '6px',
                  borderLeft: '4px solid var(--unknown)',
                  fontSize: '13px',
                }}
              >
                <div style={{ fontWeight: 600, color: 'var(--unknown)', marginBottom: '2px' }}>
                  {u.statement}
                </div>
                <div style={{ color: 'var(--text-2)' }}>
                  Reason: {u.reason}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Next Safe Action Callout */}
      <div
        style={{
          marginTop: '24px',
          padding: '16px',
          background: 'var(--surface-2)',
          border: '1px solid var(--border-strong)',
          borderRadius: '8px',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <div>
            <span style={{ fontSize: '11px', textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--text-3)' }}>
              Next Safe Action (Deterministic Engine)
            </span>
            <h3 style={{ fontSize: '16px', fontWeight: 700, color: 'var(--text-1)', marginTop: '2px' }}>
              {passport.next_safe_action.label}
            </h3>
            <p style={{ fontSize: '13px', color: 'var(--text-2)', marginTop: '4px' }}>
              {passport.next_safe_action.rationale}
            </p>
          </div>
          <button
            style={{
              padding: '8px 16px',
              backgroundColor: 'var(--accent)',
              color: '#fff',
              border: 'none',
              borderRadius: '6px',
              fontWeight: 600,
              fontSize: '13px',
            }}
          >
            Acknowledge Action
          </button>
        </div>
      </div>
    </div>
  );
};

const defaultPassport: PassportType = {
  id: 'pass-8821',
  case_id: 'case-8821',
  version: 2,
  snapshot_version: 4,
  viewer_role: 'sender',
  content_sha256: '9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08',
  created_at: new Date().toISOString(),
  confirmed_facts: [
    {
      code: 'ORDER_APPROVED',
      statement: 'PayPal Order #8NW921 was approved by buyer sandbox credentials.',
      evidence_id: 'EV-ORD-8NW921',
      tier: 1,
      observed_at: new Date(Date.now() - 180000).toISOString(),
    },
    {
      code: 'CAPTURE_COMPLETED',
      statement: 'PayPal Capture #3KK109 was executed and COMPLETED on the PayPal network.',
      evidence_id: 'EV-CAP-3KK109',
      tier: 1,
      observed_at: new Date(Date.now() - 60000).toISOString(),
    },
  ],
  uncertain_facts: [
    {
      statement: 'Recipient Bank Availability & Withdrawal Status',
      reason: 'Cross-border recipient bank settlement is not observable via PayPal checkout APIs (Default NOT_OBSERVABLE).',
    },
  ],
  next_safe_action: {
    code: 'NONE_REQUIRED',
    label: 'Original Transaction Confirmed — No Retries Needed',
    rationale: 'PayBridge Reconciler matched the original capture directly with PayPal. Blind duplicate payment successfully prevented.',
    preconditions_met: true,
    designated_actor: 'ALL',
  },
};
