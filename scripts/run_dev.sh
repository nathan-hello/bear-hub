#!/bin/bash

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

log_with_timestamp "tailwindcss -i src/public/css/tw-base.css -o src/public/css/tw-output.css -w" "log/tailwind.log" "tailwindcss"
log_with_timestamp "templ generate -path src/components --watch" "log/templ.log" "templ      "
log_with_timestamp "air" "log/air.log" "air        "  


trap 'cleanup' INT # Trap interrupt signal and call cleanup function
( trap exit SIGINT ; read -r -d '' _ </dev/tty ) # wait for Ctrl-C
exit
