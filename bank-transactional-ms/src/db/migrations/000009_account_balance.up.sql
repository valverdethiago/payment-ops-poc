CREATE TABLE IF NOT EXISTS account_balance
(
  account_balance_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
  account_uuid UUID NOT NULL,
  amount NUMERIC(15,2) NOT NULL,
  currency TEXT NOT NULL,
  date_time TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  PRIMARY KEY (account_balance_uuid),
  FOREIGN KEY (account_uuid) REFERENCES account (account_uuid)
);