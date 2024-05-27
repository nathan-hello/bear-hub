#!/bin/bash

dir=${0%/*}
cd "$dir/.."

source .env

DB_URI=$DB_URI go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate 
