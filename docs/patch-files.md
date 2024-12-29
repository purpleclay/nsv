---
icon: material/file-edit-outline
description: Automatically patch files in your repository with the next semantic version
status: new
---

# Patch files with the next semantic version

<span class="rounded-pill">:material-test-tube: experimental</span>

Let `nsv` patch files in your repository with the next calculated semantic version by executing a custom hook:

=== "ENV"

    ```{ .sh .no-select }
    NSV_HOOK="./scripts/patch.sh" nsv patch
    ```

=== "CLI"

    ```{ .sh .no-select }
    nsv patch --hook "./scripts/patch.sh"
    ```

Any file changes are committed with the default message `chore: patched files for release <version> [skip ci]`[^1].

[^1]: `nsv` detects the CI platform and changes the commit suffix accordingly.

!!! tip "Auto-patching is on the horizon."

    Soon, `nsv` will recognize standard project files and automatically patch them
    with the next semantic version. How cool is that! :material-sunglasses:

## Signing your commit

If you require GPG signing, you can configure it [here](./git-signing.md).

## Using a custom commit message

You can change the commit message. Support for Go templating provides extra [customization](./reference/templating.md#commit-message).

=== "ENV"

    ```{ .sh .no-select }
    NSV_COMMIT_MESSAGE="chore: bumped files to {{.Tag}}" nsv patch
    ```

=== "CLI"

    ```{ .sh .no-select }
    nsv patch --commit-message "chore: bumped files to {{.Tag}}"
    ```

Resulting in a commit message of:

```{ .text .no-select .no-copy }
chore: bumped files to 0.2.0
```

## Version template customization

Internally, `nsv` utilizes a go template to construct the next semantic version. Runtime customization of this template is available [here](./next-version.md#version-template-customization).
