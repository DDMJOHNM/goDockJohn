CREATE TABLE "users"(
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "created-at" timestamp NOT NULL DEFAULT (now())
);
