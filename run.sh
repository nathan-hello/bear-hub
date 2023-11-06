#!/bin/bash

# Define the commands
TAILWIND_CMD="bun run tailwindcss -o src/static/css/output.css --watch"
AIR_CMD="air"

# Run tailwindcss with prefix in the background and redirect its output to a file
$TAILWIND_CMD > tailwind.log 2>&1 &
TAILWIND_PID=$!

# Run air with prefix in the background and redirect its output to a file
$AIR_CMD > air.log 2>&1 &
AIR_PID=$!

# Function to display logs with prefixes
tail_logs_with_prefix() {
    tail -f tailwind.log -f air.log | sed -e 's/^/[Tailwind] /' -e 's/^/[Air] /'
}

# Function to kill both processes
cleanup() {
    echo "Stopping Tailwind and Air..."
    kill $TAILWIND_PID $AIR_PID
    pkill -P $$ # Kills any subprocesses started by this script
    wait $TAILWIND_PID $AIR_PID 2>/dev/null
    echo "Processes have been stopped."
    exit
}

# Trap CTRL+C and call the cleanup function
trap cleanup SIGINT

# Tail the log files with prefixes
tail_logs_with_prefix

# Wait for both processes to exit
wait $TAILWIND_PID $AIR_PID

