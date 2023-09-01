#!/bin/bash

# Globals
_file=./secrets/GITHUB_TOKEN

# Helper Functions
function msg() {
    echo "[setup-github-token.sh INFO]: $@" >&2
}

# Script
if [ -f $_file ]
then
    msg "$_file already set"
    exit 0
fi

# Read Password
echo -n 'Please enter your GitHub Token: '
read -s password
echo

# Run Command
echo $password > $_file

