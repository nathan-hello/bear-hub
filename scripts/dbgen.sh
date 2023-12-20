#!/bin/bash

dir=${0%/*}
cd "$dir/.."

source .env

pg_dump $DB_URI -s -f src/sqlc/input/remote-schema.sql
DB_URI=$DB_URI sqlc generate

