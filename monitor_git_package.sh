#!/bin/bash

echo "🔍 Monitoring git package tests specifically"
echo "==========================================="

LOG_FILE="git_package_monitor.log"
> "$LOG_FILE"

monitor_git_processes() {
    while true; do
        PROCESSES=$(ps -ax -o pid,command | grep -E '\bgit\b' | grep -v grep | grep -v monitor | grep -v bash)
        if [ ! -z "$PROCESSES" ]; then
            echo "$(date '+%H:%M:%S.%3N') - Git processes:" >> "$LOG_FILE"
            echo "$PROCESSES" >> "$LOG_FILE"
            echo "---" >> "$LOG_FILE"

            echo "⚠️  Git processes at $(date '+%H:%M:%S'):"
            echo "$PROCESSES"
            echo ""
        fi
        sleep 0.1
    done
}

echo "🚀 Starting git package test monitoring..."
monitor_git_processes &
MONITOR_PID=$!

echo "📊 Monitor PID: $MONITOR_PID"
echo ""

# Test only git package
go test ./git -v

kill $MONITOR_PID 2>/dev/null
wait $MONITOR_PID 2>/dev/null

echo ""
if [ -s "$LOG_FILE" ]; then
    echo "❌ Git processes detected in git package tests!"
    echo "============================================="
    cat "$LOG_FILE"
    echo "============================================="
else
    echo "✅ No git processes detected in git package tests"
    rm -f "$LOG_FILE"
fi
