#!/usr/bin/env bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE TABLE "user" (
      id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
      username VARCHAR(50) UNIQUE NOT NULL,
      email TEXT UNIQUE NOT NULL,
      active BOOLEAN DEFAULT true,
      created_at TIMESTAMPTZ DEFAULT NOW()
  );
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER "johndoe" WITH PASSWORD 'johndoe';
	CREATE DATABASE "johndoe";
	GRANT ALL PRIVILEGES ON DATABASE "johndoe" TO "johndoe";
EOSQL

