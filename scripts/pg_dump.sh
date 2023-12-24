#!/bin/bash

dir=${0%/*}
cd "$dir/.."

source .env
touch src/sqlc/input/remote-schema.sql

pg_dump $DB_URI --schema-only --format=plain --file src/sqlc/input/remote-schema.sql --enable-row-security
DB_URI=$DB_URI sqlc generate