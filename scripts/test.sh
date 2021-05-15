#!/bin/bash

while read module; do
  cd ./${module}
  go test -v ./...
  cd ../
done < MODULES.txt
