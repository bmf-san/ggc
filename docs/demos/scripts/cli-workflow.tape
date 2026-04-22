# ggc CLI workflow showcase
#
# Requires:
#   GGC_DEMO_CLI_WORKDIR - path to the prepared fixture workspace
#
# The goal of this tape is to contrast ggc's space-separated subcommand
# syntax against git, highlighting the shortcuts that ggc actually ships:
#   - `ggc status`
#   - `ggc diff --staged`
#   - `ggc add .`
#   - `ggc commit "<msg>"`          (no -m required)
#   - `ggc log graph`               (visual log with branches)

Output docs/demos/generated/cli-workflow.gif
Set Shell "bash"
Set Width 1200
Set Height 675
Set FontSize 22
Set TypingSpeed 60ms
Set Framerate 14
Set Theme "Dracula"

Type "cd $GGC_DEMO_CLI_WORKDIR"
Enter
Sleep 400ms

# What has changed?
Type "ggc status"
Enter
Sleep 1800ms

# Stage everything in one go
Type "ggc add ."
Enter
Sleep 1200ms

# Confirm the stage
Type "ggc status"
Enter
Sleep 1600ms

# Commit with a quoted message - note: no -m flag needed
Type "ggc commit 'docs: polish readme'"
Enter
Sleep 1600ms

# A prettier log than git log --oneline
Type "ggc log graph"
Enter
Sleep 2200ms
Type "q"
Sleep 400ms

Sleep 1200ms
