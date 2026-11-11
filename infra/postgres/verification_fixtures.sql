-- PayBridge Enterprise Verification Fixtures
-- Tests: Prerequisite foreign keys, Ledger balancing trigger (0011), and Debezium CDC outbox record.

-- 1. Create Prerequisite Tenant and Users
INSERT INTO tenants (id, name, slug)
VALUES ('11111111-1111-1111-1111-111111111111', 'Acme Enterprise Payments', 'acme-corp')
ON CONFLICT (id) DO NOTHING;

INSERT INTO users (id, email, display_name)
VALUES 
  ('22222222-2222-2222-2222-222222222221', 'sender@acmepay.com', 'Alice Sender'),
  ('22222222-2222-2222-2222-222222222222', 'receiver@acmepay.com', 'Bob Receiver')
ON CONFLICT (id) DO NOTHING;

-- 2. Create Prerequisite Payment Case
INSERT INTO payment_cases (
  id,
  tenant_id,
  external_reference,
  sender_user_id,
  receiver_user_id,
  sender_country,
  receiver_country,
  payment_purpose,
  amount,
  currency,
  environment,
  provider_state,
  application_state,
  receiver_state,
  case_status
) VALUES (
  '33333333-3333-3333-3333-333333333333',
  '11111111-1111-1111-1111-111111111111',
  'EXT-CASE-2026-001',
  '22222222-2222-2222-2222-222222222221',
  '22222222-2222-2222-2222-222222222222',
  'US',
  'IN',
  'P0802',
  5000.000000,
  'USD',
  'REAL_SANDBOX',
  'NOT_CREATED',
  'SYNCHRONISED',
  'NOT_OBSERVABLE',
  'OPEN'
) ON CONFLICT (id) DO NOTHING;

-- 3. Create Prerequisite Ledger Accounts
INSERT INTO ledger_accounts (id, tenant_id, account_number, account_name, currency, category)
VALUES
  ('44444444-4444-4444-4444-444444444441', '11111111-1111-1111-1111-111111111111', '1010-CASH-USD', 'Operating Cash Clearing', 'USD', 'ASSET'),
  ('44444444-4444-4444-4444-444444444442', '11111111-1111-1111-1111-111111111111', '2010-PAYABLE-USD', 'Customer Funds Obligation', 'USD', 'LIABILITY')
ON CONFLICT (id) DO NOTHING;

-- 4. Test Balanced Ledger Journal (MUST SUCCEED)
BEGIN;
INSERT INTO ledger_transactions (id, tenant_id, case_id, reference, description)
VALUES (
  '55555555-5555-5555-5555-555555555551',
  '11111111-1111-1111-1111-111111111111',
  '33333333-3333-3333-3333-333333333333',
  'TX-BALANCED-001',
  'Verified 2-leg customer funding settlement'
);

INSERT INTO ledger_entries (id, transaction_id, account_id, direction, amount, currency)
VALUES
  ('66666666-6666-6666-6666-666666666661', '55555555-5555-5555-5555-555555555551', '44444444-4444-4444-4444-444444444441', 'DEBIT', 5000.000000, 'USD'),
  ('66666666-6666-6666-6666-666666666662', '55555555-5555-5555-5555-555555555551', '44444444-4444-4444-4444-444444444442', 'CREDIT', 5000.000000, 'USD');
COMMIT;

-- 5. Test Debezium CDC Outbox Injection (Satisfies foreign keys)
INSERT INTO outbox (tenant_id, case_id, kind, payload)
VALUES (
  '11111111-1111-1111-1111-111111111111',
  '33333333-3333-3333-3333-333333333333',
  'RECONCILE',
  '{"event_id":"cdc-live-verified-01","case_id":"33333333-3333-3333-3333-333333333333","amount":5000,"currency":"USD"}'::jsonb
);

-- 6. Unbalanced Ledger Journal Negative Test (Verifies Deferred Trigger Rejection & Rollback)
DO $$
BEGIN
  INSERT INTO ledger_transactions (id, tenant_id, case_id, reference, description)
  VALUES (
    '55555555-5555-5555-5555-555555555552',
    '11111111-1111-1111-1111-111111111111',
    '33333333-3333-3333-3333-333333333333',
    'TX-UNBALANCED-001',
    'Intentional imbalance validation'
  );

  INSERT INTO ledger_entries (id, transaction_id, account_id, direction, amount, currency)
  VALUES (
    '66666666-6666-6666-6666-666666666663',
    '55555555-5555-5555-5555-555555555552',
    '44444444-4444-4444-4444-444444444441',
    'DEBIT',
    5000.000000,
    'USD'
  );

  -- Force deferred constraint evaluation immediately inside block
  SET CONSTRAINTS trg_verify_ledger_balance IMMEDIATE;

  RAISE EXCEPTION 'TEST FAILED: Unbalanced transaction was unexpectedly allowed!';
EXCEPTION
  WHEN OTHERS THEN
    IF SQLERRM LIKE '%is out of balance%' THEN
      RAISE NOTICE 'SUCCESS: Trigger trg_verify_ledger_balance correctly rejected unbalanced transaction: %', SQLERRM;
    ELSE
      RAISE;
    END IF;
END;
$$;
