BEGIN;

CREATE TABLE IF NOT EXISTS "transactions" (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY ,
  merchant_id uuid NOT NULL,
  reference TEXT NOT NULL,
  bank_reference TEXT DEFAULT NULL UNIQUE,
  type TEXT NOT NULL,
  amount BIGINT NOT NULL,
  status TEXT NOT NULL,
  metadata JSONB DEFAULT NULL,
  created TIMESTAMP DEFAULT NOW(),
  updated TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_transaction_merchant_id FOREIGN KEY (merchant_id) REFERENCES merchants(id)
);

COMMIT;
