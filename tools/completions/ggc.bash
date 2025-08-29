# bash completion for ggc
_ggc()
{
    local cur prev opts subopts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"


    opts="add branch clean version config hook restore diff status clean-interactive commit complete tag fetch log pull push rebase remote reset stash"

    case ${prev} in
        branch)
            subopts="current checkout checkout-remote delete delete-merged list-local list-remote"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        commit)
            subopts="allow-empty amend"
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
            subopts="files dirs"
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
            subopts="--prune"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        amend)
            subopts="--no-edit"
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
            subopts="interactive"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
    esac

    if [[ ${COMP_CWORD} == 1 ]] ; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    fi

    # Dynamic completion for branch checkout
    if [[ ${COMP_WORDS[1]} == "branch" && ${COMP_WORDS[2]} == "checkout" ]]; then
        local branches
        branches=$(ggc __complete branch 2>/dev/null)
        COMPREPLY=( $(compgen -W "${branches}" -- ${cur}) )
        return 0
    fi
    # Dynamic completion for add
    if [[ ${COMP_WORDS[1]} == "add" ]]; then
        local files
        files=$(ggc __complete files 2>/dev/null)
        COMPREPLY=( $(compgen -W "${files}" -- ${cur}) )
        return 0
    fi
}
complete -F _ggc ggc
