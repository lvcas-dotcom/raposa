#!/bin/sh
set -eu
root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
mkdir -p "$root/bin"
cd "$root/projects/open-rabbit/apps/cli"
go build -trimpath -o "$root/bin/open-rabbit" ./cmd/open-rabbit
go build -trimpath -o "$root/bin/ollama-provider" ./cmd/ollama-provider
