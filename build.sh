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
