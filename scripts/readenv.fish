#!/usr/bin/fish

while read -l line
    set -l parts (string split -m 1 "=" $line)
    if set -q parts[2]
        set -gx $parts[1] $parts[2]
        echo "Exporting $parts[1]"
    end
end < $argv[1]
