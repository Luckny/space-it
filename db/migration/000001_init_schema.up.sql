GRANT CONNECT ON DATABASE space_it TO space_it_api;
GRANT USAGE ON SCHEMA public TO space_it_api;

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "email" varchar(30) UNIQUE NOT NULL,
  "password" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

GRANT SELECT, INSERT ON users TO space_it_api;

CREATE TABLE "spaces" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "name" varchar(30) UNIQUE NOT NULL,
  "owner" uuid NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

GRANT SELECT, INSERT, UPDATE, DELETE ON spaces TO space_it_api;

CREATE TABLE "messages" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "space_id" uuid NOT NULL,
  "author" uuid NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

GRANT SELECT, INSERT ON messages TO space_it_api;

CREATE TABLE "request_log" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "method" varchar(10) NOT NULL,
  "path" varchar(100) NOT NULL,
  "user_id" uuid DEFAULT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

GRANT SELECT, INSERT ON request_log TO space_it_api;

CREATE TABLE "response_log" (
  "id" uuid NOT NULL,
  "status" int NOT NULL DEFAULT 0,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

GRANT SELECT, INSERT ON response_log TO space_it_api;

CREATE TABLE "permissions" (
  "space_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "read_permission" bool NOT NULL DEFAULT false,
  "write_permission" bool NOT NULL DEFAULT false,
  "delete_permission" bool NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("user_id", "space_id")
);

GRANT SELECT, INSERT ON permissions TO space_it_api;

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "spaces" ("owner");

CREATE INDEX ON "spaces" ("name");

CREATE INDEX ON "messages" ("space_id");

CREATE INDEX ON "messages" ("author");

CREATE INDEX ON "messages" ("space_id", "author");

CREATE INDEX ON "request_log" ("path");

CREATE INDEX ON "request_log" ("method");

CREATE INDEX ON "request_log" ("path", "method");

CREATE INDEX ON "response_log" ("status");

CREATE INDEX ON "permissions" ("updated_at");

ALTER TABLE "spaces" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");

ALTER TABLE "messages" ADD FOREIGN KEY ("space_id") REFERENCES "spaces" ("id");

ALTER TABLE "messages" ADD FOREIGN KEY ("author") REFERENCES "users" ("id");

ALTER TABLE "request_log" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "response_log" ADD FOREIGN KEY ("id") REFERENCES "request_log" ("id");

ALTER TABLE "permissions" ADD FOREIGN KEY ("space_id") REFERENCES "spaces" ("id");

ALTER TABLE "permissions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
