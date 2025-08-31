#!/bin/bash

# Precise git monitoring for macOS
echo "🔍 Precise Git Process Monitoring (macOS optimized)"
echo "=================================================="
echo ""

LOG_FILE="git_activity.log"
> "$LOG_FILE"

# Function to monitor git processes with macOS-compatible ps
monitor_git_precise() {
    while true; do
        # macOS compatible ps command
        PROCESSES=$(ps -ax -o pid,command | grep -E '\bgit\b' | grep -v grep | grep -v monitor | grep -v bash)
        if [ ! -z "$PROCESSES" ]; then
            TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S.%3N')
            echo "$TIMESTAMP - Git process detected:" >> "$LOG_FILE"
            echo "$PROCESSES" >> "$LOG_FILE"
            echo "---" >> "$LOG_FILE"

            echo "⚠️  Git process detected at $(date '+%H:%M:%S'):"
            echo "$PROCESSES"
            echo ""
        fi
        sleep 0.05
    done
}

# Start monitoring
echo "🚀 Starting precise monitoring..."
monitor_git_precise &
MONITOR_PID=$!

echo "📊 Monitor PID: $MONITOR_PID"
echo "📝 Log file: $LOG_FILE"
echo ""

# Run tests without cache
echo "🧪 Running fresh tests (no cache)..."
make test

# Stop monitoring
kill $MONITOR_PID 2>/dev/null
wait $MONITOR_PID 2>/dev/null

echo ""
echo "✅ Test execution completed"
echo ""

# Check results
if [ -s "$LOG_FILE" ]; then
    echo "❌ Git processes detected during test execution!"
    echo "=============================================="
    cat "$LOG_FILE"
    echo "=============================================="
    echo ""
    echo "🔧 These indicate tests with side effects that need fixing."
    exit 1
else
    echo "✅ SUCCESS: No git processes detected!"
    echo "🎉 All tests are properly isolated."
    rm -f "$LOG_FILE"
    exit 0
fi
