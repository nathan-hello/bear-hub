#!/bin/bash

dir=${0%/*}
cd "$dir/.."

source .env
touch src/db/input/remote-schema.sql

pg_dump $DB_URI  --schema-only  --format=plain --enable-row-security --file src/db/input/remote-schema.sql --schema=public

