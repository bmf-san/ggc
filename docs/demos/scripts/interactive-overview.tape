# ggc interactive workflow overview
#
# Requires:
#   GGC_DEMO_INTERACTIVE_WORKDIR - prepared fixture path
#
# This script illustrates launching interactive mode, refining search,
# building a workflow, and managing workflows in Workflow Mode.

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
Sleep 1500ms

# Search for status command and execute
Type "status"
Sleep 1500ms
Enter
Sleep 1500ms

# Exit back to shell
Ctrl+C
Sleep 1500ms

# Re-open interactive mode
Type "ggc"
Enter
Sleep 1500ms

# Build workflow: add -> commit -> push (using Tab in Search Mode)
Type "add"
Sleep 1500ms
Tab
Sleep 1500ms
Type "commit "
Sleep 1500ms
Tab
Sleep 1500ms
Type "push"
Sleep 1500ms
Tab
Sleep 1500ms

# Toggle to Workflow Mode (Ctrl+T)
Ctrl+T
Sleep 1500ms

# Create a new workflow with 'n'
Type "n"
Sleep 1500ms

# Navigate workflow list with Ctrl+N/P
Ctrl+P
Sleep 1500ms
Ctrl+N
Sleep 1500ms

# Delete a workflow with 'd'
Type "d"
Sleep 1500ms

# Execute a workflow with 'x'
Type "x"
Sleep 1500ms

# Exit interactive mode
Ctrl+C
Sleep 1500ms
