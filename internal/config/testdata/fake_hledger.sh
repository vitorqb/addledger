#!/bin/bash
# A fake hledger executable used during tests.

# Call to `files`
if [[ "$1" == "files" ]] && [[ "$#" == "1" ]]
then
    echo "/path/to/hledger.journal"
    exit 0
fi

echo "ERROR: UNEXPECTED COMMAND" >&2
exit 1
