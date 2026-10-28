#!/bin/sh
set -eu
root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$root/projects/raposa/apps/cli"
files=$(gofmt -l .)
if [ -n "$files" ]; then
  printf '%s\n' "$files"
  exit 1
fi
go vet ./...
go test -race -count=1 ./...
go build ./...
