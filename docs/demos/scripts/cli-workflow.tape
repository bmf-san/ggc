# ggc CLI workflow showcase
#
# Requires the following environment variables:
#   GGC_DEMO_CLI_WORKDIR - path to prepared fixture workspace

Output docs/demos/generated/cli-workflow.gif
Set Shell "bash"
Set Width 1200
Set Height 675
Set FontSize 22
Set TypingSpeed 55ms
Set Theme "Dracula"

Type "cd $GGC_DEMO_CLI_WORKDIR"
Enter
Sleep 350ms

Type "ls"
Enter
Sleep 1500ms

Type "ggc status"
Enter
Sleep 1500ms

Type "ggc add README.md"
Enter
Sleep 1500ms

Type "ggc status"
Enter
Sleep 1500ms

Type "ggc commit 'docs: polish readme'"
Enter
Sleep 1500ms

Type "ggc log simple"
Enter
Sleep 1500ms

Type ":q"
Enter
Sleep 1500ms

Type "ggc status"
Enter
Sleep 1500ms

Type "ggc version"
Enter
Sleep 1500ms

# Pause before ending
Sleep 1500ms
