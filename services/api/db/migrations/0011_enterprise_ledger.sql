-- 0011_enterprise_ledger.sql
-- +goose Up

CREATE TABLE ledger_accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  account_number TEXT NOT NULL UNIQUE,
  account_name TEXT NOT NULL,
  currency CHAR(3) NOT NULL,
  category TEXT NOT NULL CHECK (category IN ('ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ledger_transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  case_id UUID REFERENCES payment_cases(id) ON DELETE CASCADE,
  reference TEXT NOT NULL,
  description TEXT NOT NULL,
  occurred_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ledger_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  transaction_id UUID NOT NULL REFERENCES ledger_transactions(id) ON DELETE CASCADE,
  account_id UUID NOT NULL REFERENCES ledger_accounts(id),
  direction TEXT NOT NULL CHECK (direction IN ('DEBIT', 'CREDIT')),
  amount NUMERIC(20,6) NOT NULL CHECK (amount > 0),
  currency CHAR(3) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ix_ledger_entries_tx ON ledger_entries(transaction_id);

-- Enforce Zero-Sum Invariant on Ledger Transactions
CREATE OR REPLACE FUNCTION verify_ledger_transaction_balance() RETURNS TRIGGER AS $$
DECLARE
  v_balance NUMERIC(20,6);
BEGIN
  SELECT COALESCE(SUM(CASE WHEN direction = 'DEBIT' THEN amount ELSE -amount END), 0)
  INTO v_balance
  FROM ledger_entries
  WHERE transaction_id = NEW.transaction_id;

  -- Deferred check at transaction commit via constraint trigger
  IF v_balance <> 0 THEN
    RAISE EXCEPTION 'Ledger transaction % is out of balance: net sum is % (must be exactly zero)', NEW.transaction_id, v_balance;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE CONSTRAINT TRIGGER trg_verify_ledger_balance
  AFTER INSERT OR UPDATE ON ledger_entries
  DEFERRABLE INITIALLY DEFERRED
  FOR EACH ROW EXECUTE FUNCTION verify_ledger_transaction_balance();

-- +goose Down
DROP TRIGGER IF EXISTS trg_verify_ledger_balance ON ledger_entries;
DROP FUNCTION IF EXISTS verify_ledger_transaction_balance();
DROP TABLE IF EXISTS ledger_entries;
DROP TABLE IF EXISTS ledger_transactions;
DROP TABLE IF EXISTS ledger_accounts;

