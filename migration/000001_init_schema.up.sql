CREATE TABLE "users"(
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "createdat" timestamp NOT NULL DEFAULT (now())
);
