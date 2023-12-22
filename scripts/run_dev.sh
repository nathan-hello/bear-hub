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

    bash -c "$cmd" 2>&1 | while IFS= read -r line; do
        echo "$(date): $line"
    done >> "$logfile" &
    echo -e "$cmd @ PID: $!"
}

log_with_timestamp "tailwindcss -i src/static/css/tw-base.css -o src/static/css/tw-output.css -w" "log/dev/tailwind.log"
log_with_timestamp "templ generate -path src/components --watch" "log/dev/templ.log"
log_with_timestamp "air" "log/dev/air.log"


trap 'cleanup' INT # Trap interrupt signal and call cleanup function
( trap exit SIGINT ; read -r -d '' _ </dev/tty ) # wait for Ctrl-C
exit
