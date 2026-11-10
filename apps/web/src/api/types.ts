export interface Case {
  id: string;
  tenant_id: string;
  external_reference: string;
  sender_user_id: string;
  receiver_user_id?: string;
  sender_country: string;
  receiver_country: string;
  payment_purpose: string;
  amount: string;
  currency: string;
  environment: string;
  provider_state: string; // UNKNOWN, CAPTURE_PENDING, CAPTURE_COMPLETED, CAPTURE_DENIED
  last_confirmed_provider_state: string;
  application_state: string; // INCONSISTENT, API_TIMEOUT, SYNCHRONISED, RECONCILING
  receiver_state: string; // NOT_OBSERVABLE, CREDIT_PENDING, CREDITED
  case_status: string; // OBSERVING, RESOLVED, NEEDS_HUMAN_REVIEW, CLOSED_UNRESOLVED
  state_version: number;
  created_at: string;
  updated_at: string;
}

export interface CaseSummary {
  id: string;
  external_reference: string;
  amount: string;
  currency: string;
  provider_state: string;
  application_state: string;
  case_status: string;
  state_version: number;
  updated_at: string;
}

export interface Attempt {
  id: string;
  tenant_id: string;
  case_id: string;
  attempt_no: number;
  provider: string;
  environment: string;
  order_operation_request_id: string;
  capture_operation_request_id?: string;
  status: string; // ORDER_CREATING, ORDER_CREATED, ORDER_APPROVED, CAPTURE_SUBMITTED, COMPLETED, CANCELLED
  terminal: boolean;
  created_at: string;
  updated_at: string;
}

export interface EvidenceItem {
  id: string;
  tenant_id: string;
  case_id: string;
  kind: string; // PAYPAL_CAPTURE_CONFIRMATION, INVOICE_COPY, WEBHOOK_NOTIFICATION
  tier: number; // 1 to 6
  source_trust: string; // PROVIDER_VERIFIED, MERCHANT_REPORTED, SIMULATED
  environment: string;
  verified: boolean;
  untrusted_text: boolean;
  occurred_at?: string;
  observed_at: string;
  redacted_extract: any;
  payload_sha256: string;
  access_class: string;
}

export interface SafeAction {
  id: string;
  title: string;
  action_type: 'BLOCK_BLIND_RETRY' | 'RECONCILE_MERCHANT_RECORD' | 'REFRESH_PROVIDER_STATE' | 'REQUEST_HUMAN_REVIEW';
  risk_level: 'SAFE' | 'LOW' | 'HIGH';
  description: string;
  rationale: string;
  counterfactual_warning: string;
  recommended: boolean;
}

export interface GoldenPathStep {
  step: number;
  title: string;
  status: 'PENDING' | 'IN_PROGRESS' | 'COMPLETED' | 'FAILED';
  description: string;
  timestamp?: string;
  details?: Record<string, any>;
}
