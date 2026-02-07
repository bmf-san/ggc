# ggc branch management showcase
#
# Requires environment variable:
#   GGC_DEMO_BRANCH_WORKDIR - path to prepared fixture workspace

Output docs/demos/generated/branch-management.gif
Set Shell "bash"
Set Width 1200
Set Height 675
Set FontSize 22
Set TypingSpeed 55ms
Set Theme "Dracula"

Type "cd $GGC_DEMO_BRANCH_WORKDIR"
Enter
Sleep 350ms

Type "ggc branch list local"
Enter
Sleep 1500ms

Type "ggc branch create feature/docs-demo"
Enter
Sleep 1500ms

Type "ggc branch list local"
Enter
Sleep 1500ms

Type "ggc branch info"
Enter
Sleep 1500ms
Type "3"
Enter
Sleep 1500ms

Type "ggc branch checkout feature/docs-demo"
Enter
Sleep 1500ms
Type "3"
Enter
Sleep 1500ms

Sleep 1500ms
