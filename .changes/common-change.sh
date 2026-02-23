#!/bin/sh -eu

: "${CHANGE_DIR:=$PWD}"
: "${CHANGE_ISSUE_BASE_URL:=https://issues.example.com/CHANGE-}"
: "${CHANGE_AUTHOR_BASE_URL:=https://example.com/.well-known/webfinger?resource=}"

if test $# -eq 0; then
  echo "No change type provided" >&2
  exit 1
fi

if test $# -eq 1; then
  # no changes available
  exit 0
fi

printf '\n### %s\n\n' "$1"
shift 1

for P in $@; do
  if ! test -r "$P"; then
    echo "No such file: $P" >&2
    exit 2
  fi

  unset CHANGE ISSUE AUTHOR BREAKING
  . "$P"

  printf '%s ' '-'

  case "${BREAKING:-}" in
    true|True|TRUE|yes|Yes|YES|1) printf '**Breaking:** ' ;;
    *) ;;
  esac

  case "${CHANGE:-}" in
    "") printf "No change description available" ;;
    *) printf "${CHANGE}" ;;
  esac

  case "${ISSUE:-}" in
    ""|000) ;;
    *) printf " ([#%s](%s%s))" "${ISSUE}" "${CHANGE_ISSUE_BASE_URL}" "${ISSUE}" ;;
  esac

  case "${AUTHOR:-}" in
    ""|root|nobody|default|unknown|admin) ;;
    *) printf " ([%s](%s%s))" "${AUTHOR}" "${CHANGE_AUTHOR_BASE_URL}" "${AUTHOR}" ;;
  esac

  printf '\n'
done
