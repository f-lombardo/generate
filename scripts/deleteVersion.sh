#!/bin/bash

########################################################################################################################
# This script removes a tags version
########################################################################################################################

if [[ $# -ne 1 ]]; then
    echo 'This script removes a tags version' >&2
    echo 'Usage: ' >&2
    echo "$0 x.y.z" >&2
    echo "where x.y.z is the new semantic version" >&2
    exit 1
fi

git checkout master || exit 1

git tag -a $1 -m $1 && git push origin $1