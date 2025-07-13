# fish completion for ggc
function __ggc_complete_branches
    ggc __complete branch 2>/dev/null
end

function __ggc_complete_files
    ggc __complete files 2>/dev/null
end

# Main commands
complete -c ggc -f -a "add add-commit-push branch clean version diff status clean-interactive commit commit-push-interactive complete tag fetch log pull pull-rebase-push push rebase remote reset reset-clean stash stash-pull-pop"

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

# Add command with file completion
complete -c ggc -f -n "__fish_seen_subcommand_from add" -a "(__ggc_complete_files)"
