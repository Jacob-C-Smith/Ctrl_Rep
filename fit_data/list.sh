#!/bin/bash

case "$1" in
  -h | --help)
    echo "Usage: $0 [-h | --help] [-t | --tac]"
    exit 0
    ;;
  -t | --tac)
    cat='tac'
    shift # remove the processed argument
    ;;
  *)
    cat='cat'
    ;;
esac

printf "list\n" |
./key_value_db data.db |
grep -v '\$' |
grep fit |
$cat |
awk -F': ' '{print $1 ": " $2}'
