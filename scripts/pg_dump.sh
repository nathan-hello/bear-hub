#!/bin/bash

dir=${0%/*}
cd "$dir/.."

source .env
touch src/sqlc/input/remote-schema.sql

pg_dump $DB_URI  --schema-only  --format=plain --enable-row-security --file src/sqlc/input/remote-schema.sql --schema=public --schema=auth

