#!/usr/bin/env bash

for dir in *; do
  if ! [ -d "$dir" ]; then
    continue
  fi

  if ! [[ $(git -C "$dir" status --porcelain) ]]; then
    continue
  fi

  git add "$dir"
done

git commit -m "sync testdata: $(date) [skip ci]"

