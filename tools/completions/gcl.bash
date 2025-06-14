# bash completion for gcl
_gcl()
{
    local cur prev opts subopts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    opts="add branch commit push pull log fetch clean reset help stash rebase remote commit-push add-commit-push pull-rebase-push stash-pull-pop reset-clean"

    case ${prev} in
        branch)
            subopts="current checkout checkout-remote delete delete-merged"
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
            subopts="files dirs interactive"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        reset)
            subopts="clean"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        stash)
            subopts="trash"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        rebase)
            subopts="interactive"
            COMPREPLY=( $(compgen -W "${subopts}" -- ${cur}) )
            return 0
            ;;
        remote)
            subopts="list add remove set-url"
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
        branches=$(gcl __complete branch 2>/dev/null)
        COMPREPLY=( $(compgen -W "${branches}" -- ${cur}) )
        return 0
    fi
    # Dynamic completion for add
    if [[ ${COMP_WORDS[1]} == "add" ]]; then
        local files
        files=$(gcl __complete files 2>/dev/null)
        COMPREPLY=( $(compgen -W "${files}" -- ${cur}) )
        return 0
    fi
}
complete -F _gcl gcl