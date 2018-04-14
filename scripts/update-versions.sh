#!/bin/bash

cd $(dirname $0)
docker run -it --rm -v "$(pwd)/..:/hostfs" bash /hostfs/scripts/_update-versions.sh "$1"