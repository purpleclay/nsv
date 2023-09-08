---
icon: material/tag-check-outline
status: new
---

# Tag the Next Semantic Version

Let `nsv` tag your repository with the next calculated semantic version:

```{ .sh .no-select }
nsv tag
```

An annotated tag will be created with the default commit message of `chore: tagged release <version>`.

If you want to see what is happening under the hood:

=== "ENV"

    ```{ .sh .no-select }
    NSV_SHOW=true nsv tag
    ```

=== "CLI"

    ```{ .sh .no-select }
    nsv tag --show
    ```

```{ .text .no-select .no-copy }
0.1.0

HEAD
────────────────────────────────────────────────────
c6bfdda
fix: fix to the store

✓ a0a1e2b
   >>feat<<: new exciting search feature

83def28
ci: configure workflows

6c05c93
initialize repository
────────────────────────────────────────────────────
0.0.0
```

## Using a custom tag message

If you are not happy with the commit message, then feel free to change it:

=== "ENV"

    ```{ .sh .no-select }
    NSV_TAG_MESSAGE="chore: this is a custom message" nsv tag
    ```

=== "CLI"

    ```{ .sh .no-select }
    nsv tag --message "chore: this is a custom message"
    ```

## Signing your tag

If you require your tag to be signed, please ensure your git config is correct before running `nsv`. [gpg-import](https://github.com/purpleclay/gpg-import) is a tool you can easily integrate into your CI workflow and only needs a single environment variable.

## Version template customization

Internally `nsv` utilizes a go template when constructing the next semantic version. Runtime customization of this template is available [here](./next-version.md#version-template-customization).

## Committer impersonation :material-new-box:{.new-feature title="Feature added on the 8th of September 2023"}

When tagging your repository, `nsv` will identify the person associated with the commit that triggered the release and dynamically passes these to `git` through the `user.name` and `user.email` config settings.

Any of the following conditions will remove the need for impersonation:

1. The repository has the `user.name` and `user.email` settings already defined in git config.
1. The git environment variables `GIT_COMMITTER_NAME` and `GIT_COMMITTER_EMAIL` exist.
