#!/bin/bash

USAGE=$(free | grep Mem | awk '{print $3/$2 * 100}')

echo "{\"status\": 0, \"message\": \"${USAGE}% of memory currently in use\"}"
