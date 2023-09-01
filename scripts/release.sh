#!/bin/bash

#
# Helper Functions
function err() {
    echo "[release.sh ERROR]: $@" >&2
    exit 1
}
function msg() {
    echo "[release.sh INFO]: $@" >&2
}

# Check if the current branch is master
if [ "$(git rev-parse --abbrev-ref HEAD)" != "master" ]; then
  err "You must be on master to release"
fi
msg "On master :check:"

# Check if the current branch is clean
if [ -n "$(git status --porcelain)" ]; then
  err "You must have a clean working directory to release"
fi
msg "Working directory is clean :check:"

# Checks if the version bump is valid
BUMPTYPE="$1"
if [ "$BUMPTYPE" != "major" ] && [ "$BUMPTYPE" != "minor" ] && [ "$BUMPTYPE" != "patch" ]
then
    err "You must specify a valid version bump (major, minor, patch)"
fi
msg "Version bump is valid :check:"

# Gets the current version
CURRENT_VERSION=$(git describe --abbrev=0 --tags)

# Gets the new version
NEW_VERSION=$(${SEMVER} bump $BUMPTYPE $CURRENT_VERSION)

msg "Bumping version from $CURRENT_VERSION to $NEW_VERSION"

# Loads GITHUB_TOKEN for goreleaser
if [ -z "$GITHUB_TOKEN" ]
   if ! [ -f ./secrets/GITHUB_TOKEN ]
   then
       err "You must have a GITHUB_TOKEN environment variable or a secrets/GITHUB_TOKEN file"
   fi
   GITHUB_TOKEN=$(cat ./secrets/GITHUB_TOKEN)
then
export GITHUB_TOKEN

# Creates and pushes a new tag
git tag "v${NEW_VERSION}"
git push origin "v${NEW_VERSION}"
msg "Pushed tag v${NEW_VERSION} :check:"

# Creates a new release
${GORELEASER} release

