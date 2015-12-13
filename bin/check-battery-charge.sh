#!/bin/bash

CRITICAL_UNDER=30
WARNING_UNDER=50

while [[ $# > 0 ]]; do
    arg="$1"
    case $arg in
        -c|--crit)
            CRITICAL_UNDER="$2"
            shift
            ;;
        -w|--warn)
            WARNING_UNDER="$2"
            shift
            ;;
         *)
            echo "{\"status\": -1, \"message\": \"unrecognized argument ${arg}\"}"
            exit -1
            ;;
    esac
    shift
done


FILES=/sys/class/power_supply/*
TOT_POW=0

for f in $FILES
do
    if [[ "$f" = */BAT* ]]; then
        POW_LVL=$(cat ${f}/capacity)
        TOT_POW=$(($TOT_POW+$POW_LVL))
   fi
done

if [[ $TOT_POW -le $CRITICAL_UNDER ]]; then
    STATUS=2
elif [[ $TOT_POW -le $WARNING_UNDER ]]; then
    STATUS=1
else
    STATUS=0
fi

MESSAGE="total battery power is ${TOT_POW}% charged"
echo "{\"status\": ${STATUS}, \"message\": \"${MESSAGE}\"}"

