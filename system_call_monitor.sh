#!/bin/bash

echo "🔍 System Call Level Git Command Detection"
echo "=========================================="

# Note: This requires sudo privileges for dtruss
echo "⚠️  This script requires sudo privileges to trace system calls"
echo "🎯 Looking for execve() calls to git binaries"
echo ""

LOG_FILE="syscall_git_monitor.log"
> "$LOG_FILE"

echo "🚀 Starting system call monitoring (this may require password)..."

# Run tests with system call tracing
sudo dtruss -f -t execve make test 2>&1 | grep -E 'git\s+' | grep -v grep > "$LOG_FILE" &
TRACE_PID=$!

# Wait a moment for tracing to start
sleep 1

# Run the actual tests
make test

# Stop tracing
sudo kill $TRACE_PID 2>/dev/null
wait $TRACE_PID 2>/dev/null

echo ""
if [ -s "$LOG_FILE" ]; then
    echo "❌ Git execve() calls detected at system level!"
    echo "=============================================="
    cat "$LOG_FILE"
    echo "=============================================="
    echo ""
    echo "🔧 These are actual git binary executions that need investigation."
    exit 1
else
    echo "✅ SUCCESS: No git execve() calls detected at system level!"
    echo "🎉 Complete isolation confirmed - no git binaries executed."
    rm -f "$LOG_FILE"
    exit 0
fi
