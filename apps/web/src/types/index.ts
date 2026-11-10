export type ProviderState =
  | 'NOT_CREATED'
  | 'CREATED'
  | 'APPROVED'
  | 'AUTHORISED'
  | 'CAPTURE_PENDING'
  | 'CAPTURE_COMPLETED'
  | 'CAPTURE_DENIED'
  | 'REFUNDED'
  | 'PARTIALLY_REFUNDED'
  | 'UNKNOWN';

export type ApplicationState =
  | 'SYNCHRONISED'
  | 'WAITING_FOR_WEBHOOK'
  | 'API_TIMEOUT'
  | 'RECONCILING'
  | 'EVENT_GAP'
  | 'INCONSISTENT';

export type ReceiverState =
  | 'NOT_APPLICABLE'
  | 'NOT_OBSERVABLE'
  | 'REPORTED_PENDING'
  | 'VERIFIED_AVAILABLE'
  | 'WITHDRAWAL_IN_PROGRESS'
  | 'VERIFIED_BANK_CREDITED';

export type CaseStatus =
  | 'OPEN'
  | 'OBSERVING'
  | 'INVESTIGATING'
  | 'GUIDANCE_READY'
  | 'AWAITING_USER_ACTION'
  | 'VERIFYING_RESOLUTION'
  | 'RESOLVED'
  | 'NEEDS_HUMAN_REVIEW'
  | 'PAUSED_INSUFFICIENT_EVIDENCE'
  | 'CLOSED_UNRESOLVED'
  | 'ARCHIVED'
  | 'REOPENED';

export type EvidenceTier = 1 | 2 | 3 | 4 | 5 | 6;

export interface PaymentCase {
  id: string;
  tenant_id: string;
  external_reference?: string;
  sender_user_id: string;
  receiver_user_id?: string;
  sender_country: string;
  receiver_country: string;
  payment_purpose: string;
  amount: string;
  currency: string;
  environment: 'REAL_SANDBOX' | 'PAYPAL_SIMULATOR' | 'LOCAL_SIMULATION' | 'PRODUCTION';
  provider_state: ProviderState;
  last_confirmed_provider_state: ProviderState;
  provider_state_evidence_id?: string;
  provider_state_observed_at?: string;
  application_state: ApplicationState;
  receiver_state: ReceiverState;
  case_status: CaseStatus;
  state_version: number;
  created_at: string;
  updated_at: string;
}

export interface SafeAction {
  code: string;
  label: string;
  rationale: string;
  preconditions_met: boolean;
  designated_actor: string;
}

export interface SafeActionSet {
  primary_action: SafeAction;
  allowed: SafeAction[];
  blocked: SafeAction[];
}

export interface ConfirmedFact {
  code: string;
  statement: string;
  evidence_id: string;
  tier: EvidenceTier;
  observed_at: string;
}

export interface UncertainFact {
  statement: string;
  reason: string;
}

export interface PaymentTruthPassport {
  id: string;
  case_id: string;
  version: number;
  snapshot_version: number;
  viewer_role: 'sender' | 'receiver' | 'ops';
  confirmed_facts: ConfirmedFact[];
  uncertain_facts: UncertainFact[];
  next_safe_action: SafeAction;
  content_sha256: string;
  created_at: string;
}

export interface TruthGapNode {
  id: string;
  lane: 'Provider' | 'Merchant DB' | 'Browser' | 'Recipient';
  label: string;
  status: string;
  tier: EvidenceTier;
  isAmbiguous?: boolean;
}

export interface TruthGapEdge {
  from: string;
  to: string;
  relation: 'CORROBORATES' | 'CONTRADICTS' | 'PRECEDES';
}

export interface TruthGapItem {
  id: string;
  title: string;
  severity: 'ATTENTION' | 'CRITICAL';
  description: string;
}

export interface TimelineEvent {
  id: string;
  event_type: string;
  status_observed: string;
  source_trust: string;
  tier: EvidenceTier;
  verified: boolean;
  occurred_at: string;
  summary: string;
}
