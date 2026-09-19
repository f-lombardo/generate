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

gh release delete $1 --cleanup-tag
git tag -d $1
git push origin --delete $1
