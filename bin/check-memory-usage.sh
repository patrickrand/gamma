#!/bin/bash

WARNING_OVER=75
CRITICAL_OVER=90

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
        *)
            echo "{\"code\": -1, \"message\": \"unrecognized argument ${arg}\"}"
            exit -1
            ;;
    esac
    shift
done
        
USAGE=$(free | grep Mem | awk '{print $3/$2 * 100}')
function is_gte { 
    echo $(awk -v a=$1 -v b=$2 'BEGIN{print (a>b)?1:0}') 
}

if [[ $(is_gte "$USAGE" "$CRITICAL_OVER") = 1 ]]; then
    CODE=2
elif [[ $(is_gte "$USAGE" "$WARNING_OVER") = 1 ]]; then
    CODE=1
else
    CODE=0
fi

MESSAGE="${USAGE}% of memory currently in use"
echo "{\"code\": $CODE, \"message\": \"${MESSAGE}\"}"
