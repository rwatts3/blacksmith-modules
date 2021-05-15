#!/bin/bash

while read module; do
  cd ./${module}
  rm -rf go.sum
  go mod tidy
  cd ../
done < MODULES.txt
