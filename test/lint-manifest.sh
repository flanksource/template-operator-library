#!/bin/bash

make
if [ $(git status | grep "working tree clean" | wc -l) -ne "1" ]; then
  printf "Changes to manifests required, please build and recommit:\n\n"
  git status
  git diff
fi