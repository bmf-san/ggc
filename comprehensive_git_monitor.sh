#!/bin/bash

# Comprehensive git monitoring script
# Monitors processes, file system access, and network activity

echo "🔍 Comprehensive Git Process & System Call Monitoring"
echo "======================================================"
echo ""

# Create log files
PROCESS_LOG="git_processes.log"
EXEC_LOG="git_executions.log"
> "$PROCESS_LOG"
> "$EXEC_LOG"

# Function to monitor git processes with detailed info
monitor_processes() {
    while true; do
        # More detailed process monitoring
        PROCESSES=$(ps -eo pid,ppid,cmd,etime | grep -E '\bgit\b' | grep -v grep | grep -v monitor)
        if [ ! -z "$PROCESSES" ]; then
            echo "$(date '+%Y-%m-%d %H:%M:%S.%3N') - Git processes:" >> "$PROCESS_LOG"
            echo "$PROCESSES" >> "$PROCESS_LOG"
            echo "---" >> "$PROCESS_LOG"

            echo "⚠️  Git processes detected:"
            echo "$PROCESSES"
            echo ""
        fi
        sleep 0.05  # Check every 50ms for higher precision
    done
}

# Function to monitor file system for git-related activity
monitor_fs() {
    # Monitor common git directories and files
    if command -v fswatch >/dev/null 2>&1; then
        fswatch -r .git/ 2>/dev/null | while read file; do
            echo "$(date '+%Y-%m-%d %H:%M:%S.%3N') - Git file access: $file" >> "$EXEC_LOG"
            echo "📁 Git file access: $file"
        done &
        FS_PID=$!
    fi
}

# Start monitoring
echo "🚀 Starting comprehensive monitoring..."
monitor_processes &
PROCESS_PID=$!

monitor_fs
FS_PID=${FS_PID:-}

echo "📊 Process monitor PID: $PROCESS_PID"
if [ ! -z "$FS_PID" ]; then
    echo "📊 FS monitor PID: $FS_PID"
fi
echo "📝 Process log: $PROCESS_LOG"
echo "📝 Execution log: $EXEC_LOG"
echo ""

# Run tests with timing
echo "🧪 Executing tests..."
START_TIME=$(date +%s)
make test
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

# Stop monitoring
kill $PROCESS_PID 2>/dev/null
if [ ! -z "$FS_PID" ]; then
    kill $FS_PID 2>/dev/null
fi

echo ""
echo "✅ Test execution completed in ${DURATION} seconds"
echo ""

# Analyze results
HAS_ISSUES=false

if [ -s "$PROCESS_LOG" ]; then
    echo "❌ Git processes detected during test execution:"
    echo "=============================================="
    cat "$PROCESS_LOG"
    echo ""
    HAS_ISSUES=true
fi

if [ -s "$EXEC_LOG" ]; then
    echo "❌ Git file system activity detected:"
    echo "===================================="
    cat "$EXEC_LOG"
    echo ""
    HAS_ISSUES=true
fi

if [ "$HAS_ISSUES" = true ]; then
    echo "🔧 Action required: Tests have side effects and need to be fixed."
    echo "💡 Check the logs above to identify which tests are executing real git commands."
    exit 1
else
    echo "✅ SUCCESS: No git activity detected during test execution!"
    echo "🎉 All tests are properly isolated without side effects."
    rm -f "$PROCESS_LOG" "$EXEC_LOG"
    exit 0
fi
