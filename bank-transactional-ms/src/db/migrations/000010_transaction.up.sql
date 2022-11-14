CREATE TABLE IF NOT EXISTS transaction
(
  transaction_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
  account_uuid UUID NOT NULL,
  provider_account_id TEXT NOT NULL,
  description TEXT,
  description_type TEXT,
  identification TEXT NOT NULL,
  status TEXT NOT NULL,
  amount NUMERIC(15,2) NOT NULL,
  currency TEXT NOT NULL,
  date_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  PRIMARY KEY (transaction_uuid),
  FOREIGN KEY (account_uuid) REFERENCES account (account_uuid)
);