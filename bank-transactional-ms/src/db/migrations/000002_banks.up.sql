CREATE TABLE IF NOT EXISTS bank
(
  bank_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  country_code TEXT NOT NULL,
  PRIMARY KEY(bank_uuid),
  UNIQUE(code)
);