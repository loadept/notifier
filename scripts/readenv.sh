#!/usr/bin/bash

while IFS= read -r line
do
    echo "Exporting ${line%=*}"
    export $line
done < $1
