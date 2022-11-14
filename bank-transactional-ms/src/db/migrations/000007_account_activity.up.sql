CREATE TABLE IF NOT EXISTS account_activity
(
  account_activity_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
  account_uuid UUID NOT NULL,
  activity_type account_activity_type NOT NULL,
  date_time TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  PRIMARY KEY (account_activity_uuid),
  FOREIGN KEY (account_uuid) REFERENCES account (account_uuid)
);