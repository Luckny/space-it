Create Table tokens (
  "token_id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user_id" uuid NOT NULL,
  "expiry" timestamptz NOT NULL,
  "attributes" VARCHAR(4096) NOT NULL
);

GRANT SELECT, INSERT , DELETE ON tokens TO space_it_api;
