CREATE TABLE IF NOT EXISTS configuration
(
  configuration_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
  bank_uuid UUID NOT NULL,
  kafka_input_topic_name TEXT NOT NULL,
  PRIMARY KEY (configuration_uuid),
  UNIQUE(bank_uuid),
  FOREIGN KEY (bank_uuid) REFERENCES bank (bank_uuid)
);