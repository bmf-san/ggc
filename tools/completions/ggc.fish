# fish completion for ggc
function __ggc_complete_branches
    ggc __complete branch 2>/dev/null
end

function __ggc_complete_files
    ggc __complete files 2>/dev/null
end

# Main commands
complete -c ggc -f -a "add branch clean version restore hook diff status commit complete tag fetch log pull push rebase remote reset stash"

# Branch subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from branch" -a "current checkout delete rename move set info list sort contains"

# Branch checkout completion with dynamic branch names and keyword 'remote'
complete -c ggc -f -n "__fish_seen_subcommand_from branch; and __fish_seen_subcommand_from checkout" -a "remote (__ggc_complete_branches)"

# Branch list 'local', 'remote', 'verbose'
complete -c ggc -f -n "__fish_seen_subcommand_from branch; and __fish_seen_subcommand_from list" -a "local remote verbose"

# Branch delete 'merged'
complete -c ggc -f -n "__fish_seen_subcommand_from branch; and __fish_seen_subcommand_from delete" -a "merged"

# Branch set 'upstream'
complete -c ggc -f -n "__fish_seen_subcommand_from branch; and __fish_seen_subcommand_from set" -a "upstream"

# Commit subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from commit" -a "allow amend"
complete -c ggc -f -n "__fish_seen_subcommand_from commit; and __fish_seen_subcommand_from allow" -a "empty"
complete -c ggc -f -n "__fish_seen_subcommand_from commit; and __fish_seen_subcommand_from amend" -a "no-edit"

# Push subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from push" -a "current force"

# Pull subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from pull" -a "current rebase"

# Log subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from log" -a "simple graph"

# Clean subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from clean" -a "files dirs interactive"

# Complete subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from complete" -a "bash zsh fish"

# Remote subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from remote" -a "list add remove set-url"

# Fetch subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from fetch" -a "prune"

# Rebase subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from rebase" -a "interactive"

# Tag subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from tag" -a "create delete show list annotated push"

# Hook subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from hook" -a "list edit install uninstall enable disable"

# Stash subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from stash" -a "list show apply pop drop branch push save clear create store"

# Restore subcommands
complete -c ggc -f -n "__fish_seen_subcommand_from restore" -a "staged"

# Add subcommands and file completion
complete -c ggc -f -n "__fish_seen_subcommand_from add" -a "interactive patch"
complete -c ggc -f -n "__fish_seen_subcommand_from add" -a "(__ggc_complete_files)"
