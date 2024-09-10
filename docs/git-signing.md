---
icon: material/file-sign
description: Adhere to best practice with commit or tag gpg signing
social:
  cards: false
---

# Git GPG signing

If you require GPG signing, please ensure your git config is correct before running `nsv`.

## Importing a GPG key

[gpg-import](https://github.com/purpleclay/gpg-import) is a tool you can easily integrate into your CI workflow and only needs a single environment variable (`GPG_PRIVATE_KEY`) to import a GPG key and configure your git config.

## Committer impersonation

When tagging your repository, `nsv` will identify the person associated with the commit that triggered the release and dynamically passes these to `git` through the `user.name` and `user.email` config settings.

Any of the following conditions will remove the need for impersonation:

1. The repository has the `user.name` and `user.email` settings already defined in git config.
1. The git environment variables `GIT_COMMITTER_NAME` and `GIT_COMMITTER_EMAIL` exist.
