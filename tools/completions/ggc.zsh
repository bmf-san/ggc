#compdef ggc

_ggc() {
    local context state line
    typeset -A opt_args

    _arguments -C \
        '1: :_ggc_commands' \
        '*::arg:->args'

    case $state in
        args)
            case $line[1] in
                branch)
                    _ggc_branch
                    ;;
                commit)
                    _ggc_commit
                    ;;
                push)
                    _ggc_push
                    ;;
                pull)
                    _ggc_pull
                    ;;
                hook)
                    _ggc_hook
                    ;;
                status)
                    _ggc_status
                    ;;
                diff)
                    _ggc_diff
                    ;;
                log)
                    _ggc_log
                    ;;
                clean)
                    _ggc_clean
                    ;;
                complete)
                    _ggc_complete
                    ;;
                remote)
                    _ggc_remote
                    ;;
                fetch)
                    _ggc_fetch
                    ;;
                tag)
                    _ggc_tag
                    ;;
                config)
                    _ggc_config
                    ;;
                add)
                    _ggc_add
                    ;;
            esac
            ;;
    esac
}

_ggc_commands() {
    local commands
    commands=(
        'add:Add files to staging area'
        'branch:Branch operations'
        'clean:Clean working directory'
        'version:Show version information'
        'config:Configuration management'
        'hook:Git hook management'
        'diff:Show differences'
        'status:Show repository status'
        'clean-interactive:Interactive clean'
        'commit:Create commits'
        'complete:Shell completion'
        'tag:Tag management'
        'fetch:Fetch from remote'
        'log:Show commit history'
        'pull:Pull from remote'
        'push:Push to remote'
        'rebase:Rebase commits'
        'remote:Remote repository management'
        'reset:Reset changes'
        'stash:Stash changes'
    )
    _describe 'commands' commands
}

_ggc_branch() {
    local subcommands
    subcommands=(
        'current:Show current branch'
        'checkout:Switch to branch'
        'checkout-remote:Checkout remote branch'
        'delete:Delete branch'
        'delete-merged:Delete merged branches'
        'list-local:List local branches'
        'list-remote:List remote branches'
    )
    
    if [[ $CURRENT == 2 ]]; then
        _describe 'branch subcommands' subcommands
    elif [[ $words[2] == "checkout" && $CURRENT == 3 ]]; then
        # Dynamic completion for branch checkout
        local branches
        branches=(${(f)"$(ggc __complete branch 2>/dev/null)"})
        _describe 'branches' branches
    fi
}

_ggc_commit() {
    local subcommands
    subcommands=(
        'allow-empty:Allow empty commit'
        'tmp:Create temporary commit'
        'amend:Amend last commit'
    )
    
    if [[ $CURRENT == 2 ]]; then
        _describe 'commit subcommands' subcommands
    elif [[ $words[2] == "amend" && $CURRENT == 3 ]]; then
        local options
        options=(
            '--no-edit:Do not edit commit message'
        )
        _describe 'amend options' options
    fi
}

_ggc_status() {
    local options=(
        'short:Show short status'
    )
    _describe 'status options' options
}

_ggc_push() {
    local subcommands
    subcommands=(
        'current:Push current branch'
        'force:Force push'
    )
    _describe 'push subcommands' subcommands
}

_ggc_diff() {
    local subcommands
    subcommands=(
        'staged:Diff only staged changes'
        'unstaged:Diff only unstaged changes'
    )
    _describe 'diff subcommands' subcommands
}

_ggc_pull() {
    local subcommands
    subcommands=(
        'current:Pull current branch'
        'rebase:Pull with rebase'
    )
    _describe 'pull subcommands' subcommands
}

_ggc_hook() {
    local subcommands
    subcommands=(
        'list:List hooks'
        'edit:Edit hook'
        'install:Install hook'
        'uninstall:Uninstall hook'
        'enable:Enable hook'
        'disable:Disable hook'
    )
    _describe 'hook subcommands' subcommands
}

_ggc_log() {
    local subcommands
    subcommands=(
        'simple:Simple log format'
        'graph:Graph log format'
    )
    _describe 'log subcommands' subcommands
}

_ggc_clean() {
    local subcommands
    subcommands=(
        'files:Clean files'
        'dirs:Clean directories'
    )
    _describe 'clean subcommands' subcommands
}

_ggc_complete() {
    local subcommands
    subcommands=(
        'bash:Generate bash completion'
        'zsh:Generate zsh completion'
    )
    _describe 'completion shells' subcommands
}

_ggc_remote() {
    local subcommands
    subcommands=(
        'list:List remotes'
        'add:Add remote'
        'remove:Remove remote'
        'set-url:Set remote URL'
    )
    _describe 'remote subcommands' subcommands
}

_ggc_fetch() {
    local options
    options=(
        '--prune:Prune remote branches'
    )
    _describe 'fetch options' options
}

_ggc_tag() {
    local subcommands
    subcommands=(
        'create:Create tag'
        'delete:Delete tag'
        'show:Show tag'
        'list:List tags'
        'annotated:Create annotated tag'
        'push:Push tags'
    )
    _describe 'tag subcommands' subcommands
}

_ggc_config() {
    local subcommands
    subcommands=(
        'list:List configuration'
        'set:Set configuration value'
        'get:Get configuration value'
    )
    _describe 'config subcommands' subcommands
}

_ggc_add() {
    # Dynamic completion for add - get files from ggc
    local files
    files=(${(f)"$(ggc __complete files 2>/dev/null)"})
    if [[ ${#files[@]} -gt 0 ]]; then
        _describe 'files' files
    else
        _files
    fi
}

compdef _ggc ggc
