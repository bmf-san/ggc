# bash completion for ggc
_ggc()
{
    local cur prev opts subopts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"


    opts="add branch clean version config hook restore diff status commit complete tag fetch log pull push rebase remote reset stash"

    case ${prev} in
        branch)
            subopts="current checkout create delete rename move set info list sort contains"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        list)
            subopts="local remote verbose"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        commit)
            subopts="allow amend"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        push)
            subopts="current force"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        pull)
            subopts="current rebase"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        hook)
            subopts="list edit install uninstall enable disable"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        log)
            subopts="simple graph"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        clean)
            subopts="files dirs interactive"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        complete)
            subopts="bash zsh"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        remote)
            subopts="list add remove set-url"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        fetch)
            subopts="prune"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        amend)
            subopts="no-edit"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        tag)
            subopts="create delete show list annotated push"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        config)
            subopts="list set get"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        restore)
            subopts="staged"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        stash)
            subopts="list show apply pop drop branch push save clear create store"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        rebase)
            # Offer subcommands, and let specialized logic below
            # decide whether to also suggest branches.
            subopts="interactive continue abort skip"
            COMPREPLY+=( $(compgen -W "${subopts}" -- ${cur}) )
            ;;
    esac

    if [[ ${COMP_CWORD} == 1 ]] ; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    fi

    # Dynamic completion for branch checkout (local branches)
    if [[ ${COMP_WORDS[1]} == "branch" && ${COMP_WORDS[2]} == "checkout" && ${COMP_WORDS[3]} != "remote" ]]; then
        local branches
        branches=$(ggc __complete branch 2>/dev/null)
        COMPREPLY=( $(compgen -W "${branches}" -- ${cur}) )
        return 0
    fi
    # Support 'branch checkout remote' keyword
    if [[ ${COMP_WORDS[1]} == "branch" && ${COMP_WORDS[2]} == "checkout" ]]; then
        COMPREPLY+=( $(compgen -W "remote" -- ${cur}) )
        # Do not return; allow merging with branch names
    fi
    # Support 'branch delete merged'
    if [[ ${COMP_WORDS[1]} == "branch" && ${COMP_WORDS[2]} == "delete" ]]; then
        COMPREPLY=( $(compgen -W "merged" -- ${cur}) )
        return 0
    fi
    # Support 'branch set upstream'
    if [[ ${COMP_WORDS[1]} == "branch" && ${COMP_WORDS[2]} == "set" ]]; then
        COMPREPLY=( $(compgen -W "upstream" -- ${cur}) )
        return 0
    fi
    # Support 'commit allow empty'
    if [[ ${COMP_WORDS[1]} == "commit" && ${COMP_WORDS[2]} == "allow" ]]; then
        COMPREPLY=( $(compgen -W "empty" -- ${cur}) )
        return 0
    fi
    # Dynamic completion for add
    if [[ ${COMP_WORDS[1]} == "add" ]]; then
        local files
        files=$(ggc __complete files 2>/dev/null)
        COMPREPLY=( $(compgen -W "${files}" -- ${cur}) )
        return 0
    fi

    # Dynamic completion for rebase upstream (branch names)
    if [[ ${COMP_WORDS[1]} == "rebase" && ${COMP_CWORD} -eq 2 ]]; then
        # Only suggest branches when not selecting a subcommand
        case ${cur} in
            continue|abort|skip|interactive)
                ;;
            *)
                local branches
                branches=$(ggc __complete branch 2>/dev/null)
                COMPREPLY+=( $(compgen -W "${branches}" -- ${cur}) )
                return 0
                ;;
        esac
    fi
}
complete -F _ggc ggc
