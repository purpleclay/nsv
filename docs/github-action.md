---
icon: material/github
---

# Using the GitHub action

To get up and running within a GitHub workflow, include the publicly available `nsv-action` from the GitHub Actions [marketplace](https://github.com/marketplace/actions/nsv-next-semantic-version). You can find details on setting `inputs`, `outputs`, and `environment variables` in the documentation.

## Tagging a repository

If you wish to tag the repository without triggering another workflow, you must set the permissions of the job to `contents: write`.

```{.yaml linenums="1" hl_lines="10"}
name: ci
on:
  push:
    branches:
      - main
jobs:
  ci:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          fetch-depth: 0

      - name: NSV
        uses: purpleclay/nsv-action@v1
        env:
          GPG_PRIVATE_KEY: "${{ secrets.GPG_PRIVATE_KEY }}"
```

## Triggering another workflow

If you wish to trigger another workflow after `nsv` tags the repository, you must manually create a token (PAT) with the `public_repo` permission and use it during the checkout. For best security practice, use a short-lived token.

```{.yaml linenums="1" hl_lines="14"}
name: ci
on:
  push:
    branches:
      - main
jobs:
  ci:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          fetch-depth: 0
          token: "${{ secrets.TOKEN }}"

      - name: NSV
        uses: purpleclay/nsv-action@v1
        env:
          GPG_PRIVATE_KEY: "${{ secrets.GPG_PRIVATE_KEY }}"
          GPG_PASSPHRASE: "${{ secrets.GPG_PASSPHRASE }}"
```

## Capturing the next tag

You can capture the next tag without tagging the repository by setting the `next-only` input to true.

```{.yaml linenums="1" hl_lines="19"}
name: ci
on:
  push:
    branches:
      - main
jobs:
  ci:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          fetch-depth: 0

      - name: NSV
        id: nsv
        uses: purpleclay/nsv-action@v1
        with:
          next-only: true

      - name: Print Tag
        run: |
          echo "Next calculated tag: ${{ steps.nsv.outputs.nsv }}"
```
