---
icon: material/tag-check-outline
description: Automatically tag your repository with the next semantic version
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
0.2.0
┌───────────────┬──────────────────────────────────────────────────┐
│  0.2.0        │ > e0ba951                                        │
│  ↑↑           │   docs: document new exciting feature            │
│  0.1.0        │                                                  │
│               │ ✓ 2020953                                        │
│               │   >>feat<<: a new exciting feature               │
│               │                                                  │
│               │ > 709a467                                        │
│               │   ci: add github workflows                       │
└───────────────┴──────────────────────────────────────────────────┘
```

## Configurable paths for monorepo support

Monorepo support is important to the [design](./monorepos.md) of `nsv`. By adding support for context paths, multiple semantic versions can be resolved and tagged as a single operation within a repository.

```{ .sh .no-select }
nsv tag src/datastore src/notifications
```

Any version change will be printed to stdout as a comma separated list in context path order:

```{ .text .no-select .no-copy }
datastore/0.1.1,notifications/0.3.3
```

## Using a custom tag message

If you are not happy with the tag message, you can change it. Support for Go templating provides extra [customization](./reference/templating.md#tag-annotation-message).

=== "ENV"

    ```{ .sh .no-select }
    NSV_TAG_MESSAGE="chore: tagged release {{.Tag}} from {{.PrevTag}}" nsv tag
    ```

=== "CLI"

    ```{ .sh .no-select }
    nsv tag --tag-message "chore: tagged release {{.Tag}} from {{.PrevTag}}"
    ```

Resulting in a tag message of:

```{ .text .no-select .no-copy }
chore: tagged release 0.2.0 from 0.1.0
```

## Signing your tag

If you require your tag to be signed, please ensure your git config is correct before running `nsv`. [gpg-import](https://github.com/purpleclay/gpg-import) is a tool you can easily integrate into your CI workflow and only needs a single environment variable.

## Version template customization

Internally `nsv` utilizes a go template when constructing the next semantic version. Runtime customization of this template is available [here](./next-version.md#version-template-customization).

## Committer impersonation

When tagging your repository, `nsv` will identify the person associated with the commit that triggered the release and dynamically passes these to `git` through the `user.name` and `user.email` config settings.

Any of the following conditions will remove the need for impersonation:

1. The repository has the `user.name` and `user.email` settings already defined in git config.
1. The git environment variables `GIT_COMMITTER_NAME` and `GIT_COMMITTER_EMAIL` exist.

## Executing a custom hook :material-new-box:{.new-feature title="Feature added on the 2nd of September 2024"}

Before tagging your repository, `nsv` can execute a custom [hook](./hooks.md). If changes are detected, it will commit them, and then this new commit is tagged.

=== "ENV"

    ```{ .sh .no-select }
    NSV_HOOK="./scripts/patch.sh" nsv tag
    ```

=== "CLI"

    ```{ .sh .no-select }
    nsv tag --hook "./scripts/patch.sh"
    ```

It uses the default commit message of `chore: tagged release <version> [skip ci]`.

### Using a custom commit message

You can change the commit message. Support for Go templating provides extra [customization](./reference/templating.md#commit-message).

=== "ENV"

    ```{ .sh .no-select }
    NSV_COMMIT_MESSAGE="chore: bumped to {{.Tag}} {{.SkipPipelineTag}}" nsv tag
    ```

=== "CLI"

    ```{ .sh .no-select }
    nsv tag --commit-message "chore: bumped to {{.Tag}} {{.SkipPipelineTag}}"
    ```

Resulting in a commit message of:

```{ .text .no-select .no-copy }
chore: bumped to 0.2.0 [skip ci]
```

## Skip changes during a dry run :material-new-box:{.new-feature title="Feature added on the 2nd of September 2024"}

Run `nsv` within dry-run mode to skip tagging your repository and revert any changes a hook makes. This is perfect for testing.

=== "ENV"

    ```{ .sh .no-select }
    NSV_DRY_RUN="true" nsv tag
    ```

=== "CLI"

    ```{ .sh .no-select }
    nsv tag --dry-run
    ```
