#!/bin/sh
# Source this file in your .bashrc or .zshrc:
#   source /path/to/lsproj.sh
#
# This wrapper function handles the "cd" action by reading
# a temp file written by the lsproj binary.

lsproj() {
    local result_file="${TMPDIR:-/tmp}/lsproj-result-$$"

    LSPROJ_RESULT_FILE="$result_file" command lsproj "$@"
    local exit_code=$?

    if [ $exit_code -ne 0 ]; then
        rm -f "$result_file"
        return $exit_code
    fi

    if [ -f "$result_file" ]; then
        local dir
        dir="$(cat "$result_file")"
        rm -f "$result_file"
        if [ -d "$dir" ]; then
            cd "$dir" || return 1
        fi
    fi
}
