CREATE TABLE IF NOT EXISTS account
(
  account_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
  account_number TEXT NOT NULL,
  account_type ACCOUNT_TYPE NOT NULL,
  bank_uuid UUID NOT NULL,
  PRIMARY KEY (account_uuid),
  UNIQUE(account_number, bank_uuid),
  FOREIGN KEY (bank_uuid) REFERENCES bank (bank_uuid)
);