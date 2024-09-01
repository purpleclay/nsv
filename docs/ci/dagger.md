---
icon: material/sword
social:
  cards: false
---

# Using Dagger in your CI

## GitHub Action

Run [nsv](https://daggerverse.dev/mod/github.com/purpleclay/daggerverse/nsv) using the official Dagger [GitHub Action](https://github.com/dagger/dagger-for-github). The Dagger Cloud offers enhanced layer caching, which can be enabled by setting a `DAGGER_CLOUD_TOKEN` environment variable.

```{.yaml .no-select linenums="1"}
name: ci
on:
  push:
    branches:
      - main
jobs:
  nsv:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
          token: ${{ secrets.GH_NSV }} # (1)!

      - name: Tag
        uses: dagger/dagger-for-github@v6
        env:
          GPG_PRIVATE_KEY: "${{ secrets.GPG_PRIVATE_KEY }}"
          GPG_PASSPHRASE: "${{ secrets.GPG_PASSPHRASE }}"
        with:
          verb: call
          module: github.com/purpleclay/daggerverse/nsv
          args: --src . tag --show --gpg-private-key env:GPG_PRIVATE_KEY --gpg-passphrase env:GPG_PASSPHRASE
          cloud-token: ${{ secrets.DAGGER_CLOUD_TOKEN }}
```

1. A PAT token triggers another workflow after tagging the repository; this is optional.

## GitLab Template

The same Dagger experience is possible within GitLab using the Purple Clay [template](https://gitlab.com/purpleclay/templates/-/blob/main/dagger/README.md?ref_type=heads). The Dagger Cloud offers enhanced layer caching, which can be enabled by setting a `DAGGER_CLOUD_TOKEN` environment variable.

```{.yaml .no-select linenums="1" hl_lines="28-29"}
include:
  - "https://gitlab.com/purpleclay/templates/-/raw/dagger/0.11.6/dagger/Mixed.gitlab-ci.yml"

  nsv:
    extends: [.dagger]
    stage: release
    rules:
      - if: $CI_COMMIT_TAG
        when: never
      - if: $CI_PIPELINE_SOURCE == "schedule"
        when: never
      - if: $CI_MERGE_REQUEST_IID
        when: never
      - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
        when: on_success
    variables:
      GIT_DEPTH: 0
      GIT_STRATEGY: clone
      DAGGER_MODULE: "github.com/purpleclay/daggerverse/nsv"
      DAGGER_ARGS: >-
        --src .
        tag
        --show
        --paths ${WORKING_DIRECTORY}
        --gpg-private-key env:NSV_GPG_PRIVATE_KEY
        --gpg-passphrase env:NSV_GPG_PASSPHRASE
    before_script: # (1)!
      - PROJECT_URL=${CI_PROJECT_URL#"https://"}
      - git remote set-url origin "https://oauth2:${NSV_GITLAB_TOKEN}@${PROJECT_URL}.git"
```

1. To push a newly created tag, an access token with `:write_repository` permissions is required. Here, it is assigned to the `NSV_GITLAB_TOKEN` CI variable.
