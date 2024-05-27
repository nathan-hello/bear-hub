#!/bin/bash

# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail

dir=${0%/*}
cd "$dir/.."

cleanup() {
    echo -e "\nStopping tailwind, air, and templ "
    pkill -9 air
    pkill -9 templ
    pkill -9 tailwindcss
}

log_with_timestamp() {
    local cmd=$1
    local logfile=$2
    local program=$3

    bash -c "$cmd" 2>&1 | while IFS= read -r line; do
        echo "$program = $(date +"%T"): $line" | tee -a $logfile
    done &
    echo -e "$cmd @ PID: $!"
}

tw_log="${dir}/../tmp/log/tailwind.log"
air_log="${dir}/../tmp/log/air.log"
templ_log="${dir}/../tmp/log/templ.log"

touch $tw_log $air_log $templ_log

log_with_timestamp "tailwindcss -i public/css/tw-base.css -o public/css/tw-output.css -w" $tw_log "tailwindcss"
log_with_timestamp "templ generate -path src/components --watch" $templ_log "templ      "
log_with_timestamp "go run github.com/cosmtrek/air@v1.52.0" $air_log "air        "  


trap 'cleanup' INT # Trap interrupt signal and call cleanup function
( trap exit SIGINT ; read -r -d '' _ </dev/tty ) # wait for Ctrl-C
exit
