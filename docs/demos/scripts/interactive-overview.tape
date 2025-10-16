# ggc interactive workflow overview
#
# Requires:
#   GGC_DEMO_INTERACTIVE_WORKDIR - prepared fixture path
#
# This script illustrates launching interactive mode, refining search,
# building a workflow, and executing it.

Output docs/demos/generated/interactive-overview.gif
Set Shell "bash"
Set Width 1200
Set Height 675
Set FontSize 22
Set TypingSpeed 55ms
Set Theme "Dracula"

Type "cd $GGC_DEMO_INTERACTIVE_WORKDIR"
Enter
Sleep 350ms

# Launch interactive mode
Type "ggc"
Enter
Sleep 1400ms

# Search for status command and execute
Type "status"
Sleep 400ms
Enter
Sleep 1600ms

# Exit back to shell
Ctrl+C
Sleep 500ms

# Re-open interactive mode
Type "ggc"
Enter
Sleep 1400ms

# Build workflow: add -> commit -> push
Type "add"
Sleep 400ms
Tab
Sleep 400ms
Type "commit "
Sleep 400ms
Tab
Sleep 400ms
Type "push"
Sleep 400ms
Tab
Sleep 400ms

# Toggle workflow view and execute
Ctrl+T
Sleep 1200ms
Enter
Sleep 2600ms

Type "README.md"
Enter
Sleep 900ms

# Exit interactive mode
Ctrl+C
Sleep 400ms
