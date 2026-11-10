import { Case, CaseSummary, Attempt, EvidenceItem, SafeAction } from './types';

const API_BASE = '/api/v1';

// Initial Seeded Demo Cases
export const SEEDED_CASES: Case[] = [
  {
    id: 'c101-demo-4567-89ab-cdef00000001',
    tenant_id: '018f3a21-9b34-7a12-89cd-456789000000',
    external_reference: 'INV-2026-8941',
    sender_user_id: 'usr-sender-991',
    sender_country: 'US',
    receiver_country: 'DE',
    payment_purpose: 'SAAS_SUBSCRIPTION',
    amount: '149.50',
    currency: 'USD',
    environment: 'REAL_SANDBOX',
    provider_state: 'CAPTURE_COMPLETED',
    last_confirmed_provider_state: 'CAPTURE_COMPLETED',
    application_state: 'API_TIMEOUT',
    receiver_state: 'CREDIT_PENDING',
    case_status: 'NEEDS_HUMAN_REVIEW',
    state_version: 3,
    created_at: new Date(Date.now() - 15 * 60 * 1000).toISOString(),
    updated_at: new Date(Date.now() - 2 * 60 * 1000).toISOString(),
  },
  {
    id: 'c102-demo-4567-89ab-cdef00000002',
    tenant_id: '018f3a21-9b34-7a12-89cd-456789000000',
    external_reference: 'INV-2026-9022',
    sender_user_id: 'usr-sender-992',
    sender_country: 'US',
    receiver_country: 'GB',
    payment_purpose: 'HARDWARE_PURCHASE',
    amount: '580.00',
    currency: 'USD',
    environment: 'REAL_SANDBOX',
    provider_state: 'CAPTURE_COMPLETED',
    last_confirmed_provider_state: 'CAPTURE_COMPLETED',
    application_state: 'SYNCHRONISED',
    receiver_state: 'CREDITED',
    case_status: 'RESOLVED',
    state_version: 5,
    created_at: new Date(Date.now() - 4 * 3600 * 1000).toISOString(),
    updated_at: new Date(Date.now() - 3 * 3600 * 1000).toISOString(),
  },
  {
    id: 'c103-demo-4567-89ab-cdef00000003',
    tenant_id: '018f3a21-9b34-7a12-89cd-456789000000',
    external_reference: 'INV-2026-9105',
    sender_user_id: 'usr-sender-993',
    sender_country: 'US',
    receiver_country: 'FR',
    payment_purpose: 'CLOUD_SERVICES',
    amount: '1200.00',
    currency: 'USD',
    environment: 'REAL_SANDBOX',
    provider_state: 'UNKNOWN',
    last_confirmed_provider_state: 'UNKNOWN',
    application_state: 'INCONSISTENT',
    receiver_state: 'NOT_OBSERVABLE',
    case_status: 'OBSERVING',
    state_version: 1,
    created_at: new Date(Date.now() - 45 * 60 * 1000).toISOString(),
    updated_at: new Date(Date.now() - 10 * 60 * 1000).toISOString(),
  },
];

export const SEEDED_ATTEMPTS: Record<string, Attempt[]> = {
  'c101-demo-4567-89ab-cdef00000001': [
    {
      id: 'att-101-1',
      tenant_id: '018f3a21-9b34-7a12-89cd-456789000000',
      case_id: 'c101-demo-4567-89ab-cdef00000001',
      attempt_no: 1,
      provider: 'PAYPAL',
      environment: 'REAL_SANDBOX',
      order_operation_request_id: 'req-op-98172',
      capture_operation_request_id: 'cap-req-00129',
      status: 'CAPTURE_SUBMITTED',
      terminal: false,
      created_at: new Date(Date.now() - 14 * 60 * 1000).toISOString(),
      updated_at: new Date(Date.now() - 12 * 60 * 1000).toISOString(),
    },
  ],
};

export const SEEDED_EVIDENCE: Record<string, EvidenceItem[]> = {
  'c101-demo-4567-89ab-cdef00000001': [
    {
      id: 'evi-101-a',
      tenant_id: '018f3a21-9b34-7a12-89cd-456789000000',
      case_id: 'c101-demo-4567-89ab-cdef00000001',
      kind: 'PAYPAL_CAPTURE_CONFIRMATION',
      tier: 1,
      source_trust: 'PROVIDER_VERIFIED',
      environment: 'REAL_SANDBOX',
      verified: true,
      untrusted_text: false,
      occurred_at: new Date(Date.now() - 12 * 60 * 1000).toISOString(),
      observed_at: new Date(Date.now() - 11 * 60 * 1000).toISOString(),
      redacted_extract: {
        id: 'CAP-PAYPAL-782199',
        status: 'COMPLETED',
        amount: { value: '149.50', currency_code: 'USD' },
        seller_protection: { status: 'ELIGIBLE' },
      },
      payload_sha256: '750617fa580caf88c5ef332a67e101ab8872b2c8901f419c8369da9810a91122',
      access_class: 'TENANT',
    },
    {
      id: 'evi-101-b',
      tenant_id: '018f3a21-9b34-7a12-89cd-456789000000',
      case_id: 'c101-demo-4567-89ab-cdef00000001',
      kind: 'MERCHANT_CHECKOUT_LOG',
      tier: 4,
      source_trust: 'MERCHANT_REPORTED',
      environment: 'REAL_SANDBOX',
      verified: true,
      untrusted_text: false,
      occurred_at: new Date(Date.now() - 14 * 60 * 1000).toISOString(),
      observed_at: new Date(Date.now() - 14 * 60 * 1000).toISOString(),
      redacted_extract: {
        error: 'UPSTREAM_GATEWAY_TIMEOUT',
        duration_ms: 15200,
        merchant_order_ref: 'INV-2026-8941',
      },
      payload_sha256: '49e290f192b0c33a92178cd92b039488a0b0188941a80c921f009b821901cd01',
      access_class: 'TENANT',
    },
  ],
};

export const SEEDED_ACTIONS: Record<string, SafeAction[]> = {
  'c101-demo-4567-89ab-cdef00000001': [
    {
      id: 'act-1',
      title: 'Acknowledge Payment & Reconcile Order Record',
      action_type: 'RECONCILE_MERCHANT_RECORD',
      risk_level: 'SAFE',
      description: 'PayPal confirmed capture CAP-PAYPAL-782199 of $149.50 USD. Update internal merchant order status to PAID.',
      rationale: 'Provider ground truth is confirmed via authoritative GET /v2/payments/captures. Funds are already secured.',
      counterfactual_warning: 'If you do nothing or retry payment, customer will be charged twice, resulting in a dispute.',
      recommended: true,
    },
    {
      id: 'act-2',
      title: 'Enforce Blind-Retry Block',
      action_type: 'BLOCK_BLIND_RETRY',
      risk_level: 'SAFE',
      description: 'Block all incoming checkout re-attempts for INV-2026-8941.',
      rationale: 'Active non-terminal attempt already captured funds at provider.',
      counterfactual_warning: 'Allowing customer to click "Pay Again" creates a duplicate charge.',
      recommended: true,
    },
  ],
};

class ApiService {
  private useSimulation = false;
  private localCases: Case[] = [...SEEDED_CASES];
  private localAttempts: Record<string, Attempt[]> = { ...SEEDED_ATTEMPTS };
  private localEvidence: Record<string, EvidenceItem[]> = { ...SEEDED_EVIDENCE };

  constructor() {
    // Check if backend is reachable, fallback to simulation if offline
    this.checkBackendLiveness();
  }

  async checkBackendLiveness(): Promise<boolean> {
    try {
      const res = await fetch('/healthz', { signal: AbortSignal.timeout(1500) });
      if (res.ok) {
        this.useSimulation = false;
        return true;
      }
    } catch {
      // Backend offline
    }
    this.useSimulation = true;
    return false;
  }

  isSimulationMode(): boolean {
    return this.useSimulation;
  }

  setSimulationMode(val: boolean) {
    this.useSimulation = val;
  }

  async listCases(): Promise<CaseSummary[]> {
    if (this.useSimulation) {
      return this.localCases.map(c => ({
        id: c.id,
        external_reference: c.external_reference,
        amount: c.amount,
        currency: c.currency,
        provider_state: c.provider_state,
        application_state: c.application_state,
        case_status: c.case_status,
        state_version: c.state_version,
        updated_at: c.updated_at,
      }));
    }

    try {
      const res = await fetch(`${API_BASE}/cases`, {
        headers: { 'Authorization': 'Bearer demo_session_token' },
      });
      if (!res.ok) throw new Error('API request failed');
      const data = await res.json();
      return data.items || [];
    } catch {
      this.useSimulation = true;
      return this.listCases();
    }
  }

  async getCase(id: string): Promise<Case | null> {
    if (this.useSimulation) {
      const found = this.localCases.find(c => c.id === id);
      return found || null;
    }

    try {
      const res = await fetch(`${API_BASE}/cases/${id}`, {
        headers: { 'Authorization': 'Bearer demo_session_token' },
      });
      if (res.status === 404) return null;
      if (!res.ok) throw new Error('API request failed');
      return await res.json();
    } catch {
      this.useSimulation = true;
      return this.getCase(id);
    }
  }

  async createCase(input: {
    external_reference: string;
    amount_cents: number;
    currency: string;
    sender_country: string;
    receiver_country: string;
    payment_purpose: string;
  }): Promise<Case> {
    if (this.useSimulation) {
      const newCase: Case = {
        id: `c-sim-${Date.now()}`,
        tenant_id: '018f3a21-9b34-7a12-89cd-456789000000',
        external_reference: input.external_reference,
        sender_user_id: 'usr-sender-demo',
        sender_country: input.sender_country,
        receiver_country: input.receiver_country,
        payment_purpose: input.payment_purpose,
        amount: (input.amount_cents / 100).toFixed(2),
        currency: input.currency,
        environment: 'REAL_SANDBOX',
        provider_state: 'UNKNOWN',
        last_confirmed_provider_state: 'UNKNOWN',
        application_state: 'INCONSISTENT',
        receiver_state: 'NOT_OBSERVABLE',
        case_status: 'OBSERVING',
        state_version: 1,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };
      this.localCases.unshift(newCase);
      return newCase;
    }

    const res = await fetch(`${API_BASE}/cases`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer demo_session_token',
        'Idempotency-Key': `ui-${Date.now()}`,
      },
      body: JSON.stringify(input),
    });

    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.detail || 'Failed to create case');
    }
    return await res.json();
  }

  async listAttempts(caseId: string): Promise<Attempt[]> {
    if (this.useSimulation) {
      return this.localAttempts[caseId] || [];
    }

    try {
      const res = await fetch(`${API_BASE}/cases/${caseId}/attempts`, {
        headers: { 'Authorization': 'Bearer demo_session_token' },
      });
      if (!res.ok) return [];
      const data = await res.json();
      return data.items || [];
    } catch {
      return this.localAttempts[caseId] || [];
    }
  }

  async listEvidence(caseId: string): Promise<EvidenceItem[]> {
    if (this.useSimulation) {
      return this.localEvidence[caseId] || [];
    }

    try {
      const res = await fetch(`${API_BASE}/cases/${caseId}/evidence`, {
        headers: { 'Authorization': 'Bearer demo_session_token' },
      });
      if (!res.ok) return [];
      const data = await res.json();
      return data.items || [];
    } catch {
      return this.localEvidence[caseId] || [];
    }
  }

  async getSafeActions(caseId: string): Promise<SafeAction[]> {
    return SEEDED_ACTIONS[caseId] || [
      {
        id: 'act-default',
        title: 'Query Authoritative Provider Capture Status',
        action_type: 'REFRESH_PROVIDER_STATE',
        risk_level: 'SAFE',
        description: 'Read provider transaction ground truth directly via GET /v2/payments/captures without mutating payment.',
        rationale: 'Safe read-only inquiry to determine actual financial state.',
        counterfactual_warning: 'Do not issue a refund or recharge until capture truth is established.',
        recommended: true,
      },
    ];
  }

  // Update local case state during interactive demo walk
  updateCaseState(caseId: string, updates: Partial<Case>) {
    const idx = this.localCases.findIndex(c => c.id === caseId);
    if (idx !== -1) {
      this.localCases[idx] = {
        ...this.localCases[idx],
        ...updates,
        updated_at: new Date().toISOString(),
      };
    }
  }

  addAttempt(caseId: string, attempt: Attempt) {
    if (!this.localAttempts[caseId]) this.localAttempts[caseId] = [];
    this.localAttempts[caseId].push(attempt);
  }

  addEvidence(caseId: string, item: EvidenceItem) {
    if (!this.localEvidence[caseId]) this.localEvidence[caseId] = [];
    this.localEvidence[caseId].unshift(item);
  }
}

export const api = new ApiService();
