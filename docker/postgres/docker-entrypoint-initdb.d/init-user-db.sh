#!/usr/bin/env bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER "johndoe" WITH PASSWORD 'johndoe';
	CREATE DATABASE "test";
	GRANT ALL PRIVILEGES ON DATABASE "test" TO "johndoe";
EOSQL