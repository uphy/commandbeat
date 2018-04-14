#!/bin/sh

if [ $# != 1 ]; then
  echo Specify a version.
  exit 1
fi

if [ $(echo "$1" | grep -E "^\d+\.\d+\.\d+(\-\w+)?$" > /dev/null; echo $?) == 1 ]; then
  echo Invalid version format.
  exit 1
fi

cd $(dirname $0)/..
sed -i -e "s/^ARG VERSION=.*$/ARG VERSION=$1/" docker/Dockerfile
sed -i -e "s@uphy/commandbeat:.*@uphy/commandbeat:$1@" docker/docker-compose.yml
sed -i -e "s/var Version = \".*\"/var Version = \"$1\"/" cmd/root.go
sed -i -e "s/version: \".*\"/version: \"$1\"/" docs/_config.yml