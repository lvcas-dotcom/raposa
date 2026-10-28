#!/bin/sh
set -eu
root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
mkdir -p "$root/bin"
cd "$root/projects/raposa/apps/cli"
go build -trimpath -o "$root/bin/raposa" ./cmd/raposa
go build -trimpath -o "$root/bin/ollama-provider" ./cmd/ollama-provider
