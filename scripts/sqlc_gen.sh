#!/bin/bash

dir=${0%/*}
cd "$dir/.."

source .env

PG_DB_URI=$PG_DB_URI SQLITE_DB_URI=$SQLITE_DB_URI sqlc generate 
