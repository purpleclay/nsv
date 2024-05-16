#!/bin/sh

# Borrowed from: https://raw.githubusercontent.com/goreleaser/goreleaser/main/scripts/completions.sh
set -e
rm -rf completions
mkdir completions

# Generate the shell completion scripts
for SH in bash zsh fish; do
	go run . completion "${SH}" > "completions/nsv.${SH}"
done
