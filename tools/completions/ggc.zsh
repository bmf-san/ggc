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
                restore)
                    _ggc_restore
                    ;;
                log)
                    _ggc_log
                    ;;
                clean)
                    _ggc_clean
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
                stash)
                    _ggc_stash
                    ;;
                rebase)
                    _ggc_rebase
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
        'commit:Create commits'
        'tag:Tag management'
        'fetch:Fetch from remote'
        'log:Show commit history'
        'pull:Pull from remote'
        'push:Push to remote'
        'rebase:Rebase commits'
        'remote:Remote repository management'
        'reset:Reset changes'
        'stash:Stash changes'
		'restore:Restore working tree files'
    )
    _describe 'commands' commands
}

_ggc_branch() {
    local subcommands
    subcommands=(
        'current:Show current branch'
        'checkout:Switch to branch or remote'
        'delete:Delete branch or merged'
        'rename:Rename a branch'
        'move:Move branch to specified commit'
        'set:Set upstream'
        'info:Show detailed branch information'
        'list:Show detailed branch listing'
        'sort:List branches sorted by date or name'
        'contains:Show branches containing a commit'
    )

    if [[ $CURRENT == 2 ]]; then
        _describe 'branch subcommands' subcommands
    elif [[ $words[2] == "checkout" && $CURRENT == 3 ]]; then
        # Dynamic completion for branch checkout
        local branches
        branches=(${(f)"$(ggc __complete branch 2>/dev/null)"})
        if [[ ${#branches[@]} -gt 0 ]]; then
            _describe 'branches' branches
        fi
        _values 'keyword' remote
    elif [[ $words[2] == "delete" && $CURRENT == 3 ]]; then
        _values 'keyword' merged
    elif [[ $words[2] == "set" && $CURRENT == 3 ]]; then
        _values 'keyword' upstream
    elif [[ $words[2] == "list" && $CURRENT == 3 ]]; then
        _values 'keyword' local remote verbose
    fi
}

_ggc_commit() {
    local subcommands
    subcommands=(
        'allow:Allow empty commit'
        'amend:Amend previous commit'
    )
    if [[ $CURRENT == 3 && $words[2] == "amend" ]]; then
        _values 'keyword' no-edit
        return
    elif [[ $CURRENT == 3 && $words[2] == "allow" ]]; then
        _values 'keyword' empty
        return
    fi
    _describe 'commit subcommands' subcommands
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
        'interactive:Interactive clean'
    )
    _describe 'clean subcommands' subcommands
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
    local subcommands
    subcommands=(
        'prune:Prune remote branches'
    )
    _describe 'fetch subcommands' subcommands
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

_ggc_restore() {
    local subcommands
    subcommands=(
        'staged:Unstage file(s) (restore from HEAD to index)'
    )
    _describe 'restore subcommands' subcommands
}

_ggc_add() {
    # Subcommands for add
    local subcommands
    subcommands=(
        'interactive:Interactive add'
        'patch:Patch mode'
    )
    _describe 'add subcommands' subcommands

    # Dynamic completion for add - get files from ggc
    local files
    files=(${(f)"$(ggc __complete files 2>/dev/null)"})
    if [[ ${#files[@]} -gt 0 ]]; then
        _describe 'files' files
    else
        _files
    fi
}

_ggc_stash() {
    local subcommands
    subcommands=(
        'list:List all stashes'
        'show:Show changes in stash'
        'apply:Apply stash without removing it'
        'pop:Apply and remove stash'
        'drop:Remove stash'
        'branch:Create branch from stash'
        'push:Save changes to new stash'
        'save:Save changes to new stash with message'
        'clear:Remove all stashes'
        'create:Create stash and return object name'
        'store:Store stash object'
    )
    _describe 'stash subcommands' subcommands
}

_ggc_rebase() {
    local subcommands
    subcommands=(
        'interactive:Interactive rebase'
        'continue:Continue an in-progress rebase'
        'abort:Abort an in-progress rebase'
        'skip:Skip current patch and continue'
    )
    if [[ $CURRENT == 2 ]]; then
        # Show subcommands; also suggest branches unless the current word
        # exactly matches a known subcommand.
        _describe 'rebase subcommands' subcommands
        case $words[$CURRENT] in
            (continue|abort|skip|interactive)
                ;;
            (*)
                local branches
                branches=(${(f)"$(ggc __complete branch 2>/dev/null)"})
                if [[ ${#branches[@]} -gt 0 ]]; then
                    _describe 'branches' branches
                fi
                ;;
        esac
        return
    fi
}

compdef _ggc ggc
