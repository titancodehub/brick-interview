BEGIN;

CREATE TABLE IF NOT EXISTS ledger_entries (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  merchant_id UUID NOT NULL,
  transaction_id UUID NOT NULL,
  credit BIGINT NOT NULL,
  debit BIGINT NOT NULL,
  created TIMESTAMP DEFAULT NOW(),
  updated TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_entries_merchant_id FOREIGN KEY (merchant_id) REFERENCES merchants(id)
);

COMMIT;