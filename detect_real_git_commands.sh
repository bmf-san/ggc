#!/bin/bash

echo "🔍 Detecting Real Git Commands (excluding test processes)"
echo "========================================================"

LOG_FILE="real_git_commands.log"
> "$LOG_FILE"

monitor_real_git_commands() {
    while true; do
        # Look for actual git commands, excluding test processes
        PROCESSES=$(ps -ax -o pid,command | grep -E '\bgit\s+(status|add|commit|push|pull|branch|checkout|config|describe|rev-parse|log|diff|stash|tag|remote|fetch|clean|reset|restore)' | grep -v grep | grep -v test)
        if [ ! -z "$PROCESSES" ]; then
            echo "$(date '+%H:%M:%S.%3N') - Real git command detected:" >> "$LOG_FILE"
            echo "$PROCESSES" >> "$LOG_FILE"
            echo "---" >> "$LOG_FILE"

            echo "🚨 REAL GIT COMMAND DETECTED at $(date '+%H:%M:%S'):"
            echo "$PROCESSES"
            echo ""
        fi
        sleep 0.05
    done
}

echo "🚀 Starting real git command detection..."
monitor_real_git_commands &
MONITOR_PID=$!

echo "📊 Monitor PID: $MONITOR_PID"
echo "🎯 Looking for actual git commands like: git status, git add, git commit, etc."
echo "❌ Ignoring: test processes, git.test binaries, go test commands"
echo ""

# Run all tests
make test

kill $MONITOR_PID 2>/dev/null
wait $MONITOR_PID 2>/dev/null

echo ""
if [ -s "$LOG_FILE" ]; then
    echo "❌ REAL GIT COMMANDS DETECTED!"
    echo "============================="
    cat "$LOG_FILE"
    echo "============================="
    echo ""
    echo "🔧 These are actual git commands that need to be mocked!"
    exit 1
else
    echo "✅ SUCCESS: No real git commands detected!"
    echo "🎉 All git operations are properly mocked."
    rm -f "$LOG_FILE"
    exit 0
fi
