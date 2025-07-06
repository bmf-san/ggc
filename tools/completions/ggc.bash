# bash completion for ggc
_ggc()
{
    local cur prev opts subopts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    opts="add add-commit-push branch clean diff status clean-interactive commit commit-push-interactive complete fetch log pull pull-rebase-push push rebase remote reset reset-clean stash stash-pull-pop"

    case ${prev} in
        branch)
            subopts="current checkout checkout-remote delete delete-merged list-local list-remote"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        commit)
            subopts="allow-empty tmp"
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
