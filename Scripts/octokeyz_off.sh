#!/bin/bash

set -e

if [[ ! -f "/media/fat/MiSTer" ]]; then
    echo "This script must be run on a MiSTer system."
    exit 1
fi

if mount | grep "on / .*[(,]ro[,$]" > /dev/null; then
    mount / -o remount,rw
fi

/media/fat/Scripts/.octokeyz/octokeyz-mister -init -stop

echo "octokeyz is off and inactive at startup."
exit 0
