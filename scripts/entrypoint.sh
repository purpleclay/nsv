#!/bin/sh

if [ -n "$GPG_PRIVATE_KEY" ]; then
  gpg-import
  echo
fi

exec nsv "$@"
