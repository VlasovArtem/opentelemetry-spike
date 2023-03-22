#!/bin/bash

for i in {1..5} ; do
  if [ "$i" == "$1" ]; then
    >&2 echo "Error from container $i"
    exit 1
  else
    echo "Hello from container $i"
  fi
  sleep 1
done