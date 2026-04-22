# ggc branch management showcase
#
# Requires:
#   GGC_DEMO_BRANCH_WORKDIR - path to the prepared fixture workspace
#
# Shows ggc's flagless verb-driven branch API:
#   - `ggc branch current`      quick "which branch am I on?"
#   - `ggc branch list local`   local branches with current marked
#   - `ggc branch create <n>`   create and switch in one call
#   - `ggc branch checkout`     lists + prompts (no flag soup)
#   - `ggc branch rename`       rename the current branch

Output docs/demos/generated/branch-management.gif
Set Shell "bash"
Set Width 1200
Set Height 675
Set FontSize 22
Set TypingSpeed 60ms
Set Framerate 14
Set Theme "Dracula"

Type "cd $GGC_DEMO_BRANCH_WORKDIR"
Enter
Sleep 400ms

Type "ggc branch current"
Enter
Sleep 1200ms

Type "ggc branch list local"
Enter
Sleep 1600ms

# create-and-switch in one verb
Type "ggc branch create feature/docs-demo"
Enter
Sleep 1400ms

Type "ggc branch current"
Enter
Sleep 1200ms

# numbered selector - no need to remember the exact branch name
Type "ggc branch checkout"
Enter
Sleep 1500ms
Type "1"
Enter
Sleep 1400ms

Type "ggc branch list local"
Enter
Sleep 1800ms

Sleep 1000ms
