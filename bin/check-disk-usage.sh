#!/bin/bash

WARNING_OVER=50
CRITICAL_OVER=70
DIR=/
while [[ $# > 0 ]]
do
    arg="$1"
    case $arg in 
        -w|--warn)
            WARNING_OVER="$2"
            shift
            ;;
        -c|--crit)
            CRITICAL_OVER="$2"
            shift
            ;;
        -*)
            echo "{\"status\": -1, \"message\": \"unrecognized argument ${arg}\"}"
            exit -1
            ;;
         *)
             DIR="$arg"
    esac
    shift
done
      
USAGE=$(df ${DIR} | awk '{if (NR == 2) print $5}')
USAGE="${USAGE%?}"
function is_gte { 
    echo $(awk -v a=$1 -v b=$2 'BEGIN{print (a>b)?1:0}') 
}

if [[ $(is_gte "$USAGE" "$CRITICAL_OVER") = 1 ]]; then
    STATUS=2
elif [[ $(is_gte "$USAGE" "$WARNING_OVER") = 1 ]]; then
    STATUS=1
else
    STATUS=0
fi

MESSAGE="${USAGE}% of disk mounted at ${DIR} currently in use"
echo "{\"status\": $STATUS, \"message\": \"${MESSAGE}\"}"
