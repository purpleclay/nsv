#!/bin/sh

if [ -n "$GPG_PRIVATE_KEY" ]; then
  if ! gpg-import; then
    echo
    echo "failed to import gpg private key. exiting..."
    exit 1
  fi
  echo
fi

exec nsv "$@"
