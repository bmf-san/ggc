#!/bin/bash

# Git process monitoring script for test execution
# This script monitors for any git processes during test execution

echo "🔍 Starting git process monitoring..."
echo "📋 This will detect any real git commands executed during tests"
echo ""

# Create a log file for git processes
LOG_FILE="git_processes.log"
> "$LOG_FILE"

# Function to monitor git processes
monitor_git() {
    while true; do
        # Check for git processes and log them with timestamp
        PROCESSES=$(ps aux | grep -E '\bgit\b' | grep -v grep | grep -v monitor_git)
        if [ ! -z "$PROCESSES" ]; then
            echo "$(date '+%Y-%m-%d %H:%M:%S.%3N') - Git process detected:" >> "$LOG_FILE"
            echo "$PROCESSES" >> "$LOG_FILE"
            echo "---" >> "$LOG_FILE"

            # Also print to console for real-time monitoring
            echo "⚠️  Git process detected at $(date '+%H:%M:%S'):"
            echo "$PROCESSES"
            echo ""
        fi
        sleep 0.1  # Check every 100ms
    done
}

# Start monitoring in background
monitor_git &
MONITOR_PID=$!

echo "🚀 Starting test execution with git process monitoring..."
echo "📊 Monitor PID: $MONITOR_PID"
echo "📝 Log file: $LOG_FILE"
echo ""

# Run the tests
make test

# Stop monitoring
kill $MONITOR_PID 2>/dev/null

echo ""
echo "✅ Test execution completed"
echo ""

# Check results
if [ -s "$LOG_FILE" ]; then
    echo "❌ WARNING: Git processes were detected during test execution!"
    echo "📄 Git process log:"
    echo "===================="
    cat "$LOG_FILE"
    echo "===================="
    echo ""
    echo "🔧 These tests have side effects and need to be fixed."
    exit 1
else
    echo "✅ SUCCESS: No git processes detected during test execution!"
    echo "🎉 All tests are properly isolated without side effects."
    rm -f "$LOG_FILE"
    exit 0
fi
