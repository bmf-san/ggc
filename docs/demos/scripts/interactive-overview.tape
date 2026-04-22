# ggc interactive workflow overview
#
# Requires:
#   GGC_DEMO_INTERACTIVE_WORKDIR - prepared fixture path
#
# Show off the real differentiators of interactive mode:
#   1. Type a few letters, fuzzy-search every ggc command, hit Enter.
#   2. Tab to queue commands into a workflow from search mode.
#   3. Ctrl+T to switch to Workflow mode and execute the queue.

Output docs/demos/generated/interactive-overview.gif
Set Shell "bash"
Set Width 1000
Set Height 560
Set FontSize 20
Set TypingSpeed 55ms
Set Framerate 8
Set PlaybackSpeed 1.25
Set Theme "Dracula"

Type "cd $GGC_DEMO_INTERACTIVE_WORKDIR"
Enter
Sleep 350ms

# --- 1. Fuzzy search ------------------------------------------------------
Type "ggc"
Enter
Sleep 1000ms

# "brl" fuzzy-matches "branch list"
Type "brl"
Sleep 1000ms
Enter
Sleep 1400ms

# exit back to shell
Ctrl+C
Sleep 500ms

# --- 2. Build a workflow from search mode --------------------------------
Type "ggc"
Enter
Sleep 900ms

# queue "status"
Type "status"
Sleep 700ms
Tab
Sleep 500ms
Ctrl+U

# queue "add"
Type "add"
Sleep 700ms
Tab
Sleep 500ms
Ctrl+U

# queue "commit"
Type "commit"
Sleep 700ms
Tab
Sleep 700ms

# --- 3. Switch to Workflow mode and run ----------------------------------
Ctrl+T
Sleep 1400ms

# exit
Ctrl+C
Sleep 500ms
