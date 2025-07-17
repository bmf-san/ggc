# fish completion for ggc
function __ggc_complete_branches
    ggc __complete branch 2>/dev/null
end

function __ggc_complete_files
    ggc __complete files 2>/dev/null
end

# Main commands
complete -c ggc -f -a "add branch clean version restore hook diff status clean-interactive commit complete tag fetch log pull push rebase remote reset stash"

# Branch subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from branch" -a "current checkout checkout-remote delete delete-merged list-local list-remote"

# Branch checkout completion with dynamic branch names
complete -c ggc -f -n "__fish_seen_subcommand_from branch; and __fish_seen_subcommand_from checkout" -a "(__ggc_complete_branches)"

# Commit subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from commit" -a "allow-empty tmp amend"

# Push subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from push" -a "current force"

# Pull subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from pull" -a "current rebase"

# Log subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from log" -a "simple graph"

# Clean subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from clean" -a "files dirs"

# Complete subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from complete" -a "bash zsh fish"

# Remote subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from remote" -a "list add remove set-url"

# Fetch options
complete -c ggc -f -n "__fish_seen_subcommand_from fetch" -l prune

# Amend options
complete -c ggc -f -n "__fish_seen_subcommand_from amend" -l no-edit

# Tag subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from tag" -a "create delete show list annotated push"

# Hook subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from hook" -a "list edit install uninstall enable disable"

# Restore subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from restore" -a "staged"

# Add command with file completion
complete -c ggc -f -n "__fish_seen_subcommand_from add" -a "(__ggc_complete_files)"
