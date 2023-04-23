#!/bin/bash
if ! [ -f .env ]
then
    echo "Creating .env file..."
    cat <<EOF >.env
export ADDLEDGER_DESTFILE=$(pwd)/out/destfile
EOF
    echo "Done!"
fi

