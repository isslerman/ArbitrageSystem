#!/bin/bash

# Specify the file to monitor
file_to_monitor="log.txt"

# Use tail -f to continuously read the file
# Pipe the output to grep to search for the word "good"
# If "good" is found, play a beep sound using echo -en "\007"
tail -f "$file_to_monitor" | while read line; do
    echo "$line"
    
    third_value=$(echo "$line" | awk -F, '{print $3}')
    if (( $(echo "$third_value > 0.40" | bc -l) )); then
        echo -en "\007"
    fi
done
