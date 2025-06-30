#!/bin/bash

set -e

if [[ ! -f "/media/fat/MiSTer" ]]; then
    echo "This script must be run on a MiSTer system."
    exit 1
fi

if mount | grep "on / .*[(,]ro[,$]" > /dev/null; then
    mount / -o remount,rw
fi

/media/fat/Scripts/.octokeyz/octokeyz-mister -init &

echo "octokeyz is on and active at startup."
exit 0
