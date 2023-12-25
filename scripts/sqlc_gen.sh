#!/bin/bash

dir=${0%/*}
cd "$dir/.."

source .env

DB_URI=$DB_URI sqlc generate 
