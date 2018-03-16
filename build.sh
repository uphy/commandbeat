#!/bin/bash

export GOX_FLAGS="-arch amd64"
make crosscompile

for file in $(find build/bin -type f)
do
  name=$(basename $file)
  extension=$(echo $name | sed -e 's/[^.]*//')

  mkdir -p build/package/$name
  cp $file build/package/$name/commandbeat$extension
  cp commandbeat.yml commandbeat.reference.yml build/package/$name/
  
  mkdir -p build/dist
  tar zcf build/dist/$name.tar.gz -C build/package $name
done

# Tags locally for deploying on all commit.
VERSION=$(grep "var Version" cmd/root.go | sed -e 's/^.*"\(.*\)"$/\1/')
git config --local user.name "Yuhi Ishikura"
git config --local user.email "yuhi.ishikura@uphy.jp"
git tag -d "$VERSION"
git tag -a "$VERSION" -m "$VERSION"
