import React from 'react';
import { TimelineEvent } from '../../types';

interface LifecycleTimelineProps {
  events?: TimelineEvent[];
}

export const LifecycleTimeline: React.FC<LifecycleTimelineProps> = ({
  events = defaultTimelineEvents,
}) => {
  return (
    <div className="card" style={{ background: 'var(--surface-1)' }}>
      <div style={{ marginBottom: '16px' }}>
        <h3 style={{ fontSize: '18px', fontWeight: 600, color: 'var(--text-1)' }}>
          Payment Lifecycle Timeline
        </h3>
        <p style={{ fontSize: '13px', color: 'var(--text-2)' }}>
          Chronological causal event stream with source trust indicators.
        </p>
      </div>

      <div style={{ position: 'relative', paddingLeft: '24px' }}>
        {/* Timeline vertical spine */}
        <div
          style={{
            position: 'absolute',
            left: '7px',
            top: '8px',
            bottom: '8px',
            width: '2px',
            backgroundColor: 'var(--border)',
          }}
        />

        {events.map((evt, idx) => {
          const isTimeout = evt.event_type.includes('TIMEOUT');
          const isCompleted = evt.status_observed === 'COMPLETED';

          return (
            <div
              key={evt.id || idx}
              style={{
                position: 'relative',
                marginBottom: '20px',
                paddingLeft: '12px',
              }}
            >
              {/* Event node dot */}
              <div
                style={{
                  position: 'absolute',
                  left: '-21px',
                  top: '4px',
                  width: '12px',
                  height: '12px',
                  borderRadius: '50%',
                  backgroundColor: isCompleted
                    ? 'var(--verified)'
                    : isTimeout
                    ? 'var(--unknown)'
                    : 'var(--accent)',
                  border: '2px solid var(--surface-1)',
                }}
              />

              <div
                style={{
                  background: 'var(--surface-2)',
                  border: '1px solid var(--border)',
                  borderRadius: '6px',
                  padding: '12px 16px',
                }}
              >
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '6px' }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                    <span
                      className={`badge ${
                        isCompleted
                          ? 'badge-verified'
                          : isTimeout
                          ? 'badge-unknown'
                          : 'badge-pending'
                      }`}
                    >
                      {evt.event_type}
                    </span>
                    <span style={{ fontSize: '12px', color: 'var(--text-3)' }}>
                      Tier E{evt.tier} • {evt.source_trust}
                    </span>
                  </div>
                  <time style={{ fontSize: '11px', color: 'var(--text-3)' }}>
                    {new Date(evt.occurred_at).toLocaleTimeString()}
                  </time>
                </div>
                <p style={{ fontSize: '13px', color: 'var(--text-1)', margin: 0 }}>
                  {evt.summary}
                </p>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

const defaultTimelineEvents: TimelineEvent[] = [
  {
    id: 'evt-1',
    event_type: 'ORDER_CREATED',
    status_observed: 'CREATED',
    source_trust: 'PROVIDER_VERIFIED',
    tier: 1,
    verified: true,
    occurred_at: new Date(Date.now() - 300000).toISOString(),
    summary: 'PayPal Order #8NW921 created for $250.00 USD with intent CAPTURE.',
  },
  {
    id: 'evt-2',
    event_type: 'BUYER_APPROVED',
    status_observed: 'APPROVED',
    source_trust: 'PROVIDER_VERIFIED',
    tier: 1,
    verified: true,
    occurred_at: new Date(Date.now() - 180000).toISOString(),
    summary: 'Buyer authenticated and approved order on PayPal sandbox portal.',
  },
  {
    id: 'evt-3',
    event_type: 'CAPTURE_CALL_TIMEOUT',
    status_observed: 'API_TIMEOUT',
    source_trust: 'LOCAL_OBSERVED',
    tier: 3,
    verified: false,
    occurred_at: new Date(Date.now() - 120000).toISOString(),
    summary: 'Client received HTTP 504 Gateway Timeout during capture execution. Outcome Unknown.',
  },
  {
    id: 'evt-4',
    event_type: 'READONLY_RECONCILIATION_MATCH',
    status_observed: 'COMPLETED',
    source_trust: 'PROVIDER_VERIFIED',
    tier: 1,
    verified: true,
    occurred_at: new Date(Date.now() - 60000).toISOString(),
    summary: 'Reconciler queried GET /v2/payments/captures/3KK109; confirmed capture was completed.',
  },
];
