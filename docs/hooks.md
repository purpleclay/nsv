---
icon: material/hook
description: Extend your release workflow by patching files with the next semantic version
status: new
---

# Executing a custom hook

<span class="rounded-pill">:material-test-tube: experimental</span>

To further streamline the release workflow, `nsv` adds support for executing a custom hook before tagging a repository with the next semantic version. The hook's contents are passed to a shell interpreter, which, therefore, supports both inline shell commands or a path to a script, maximizing reconfigurability.

## A hooks context

`nsv` injects the following environment variables into the hooks context:

| Variable Name           | Description                                                        | Example  |
| ----------------------- | ------------------------------------------------------------------ | -------- |
| `NSV_PREV_TAG`          | The previous semantic version                                      | `0.2.0`  |
| `NSV_NEXT_TAG`          | The next calculated semantic version based on the commit history   | `0.2.1`  |
| `NSV_WORKING_DIRECTORY` | The current execution path. Will be a sub-directory for a monorepo | `src/ui` |

Let's patch the semantic version within a `Cargo.toml` file to put it into practice.

1. Add a script to your project (`scripts/patch.sh`) for patching the file:
  ```{ .sh .no-select }
  #!/bin/sh
  sed -i '' "s/^version = \"$NSV_PREV_TAG\"/version = \"$NSV_NEXT_TAG\"/" \
    "$NSV_WORKING_DIRECTORY/Cargo.toml"
  ```

1. Run `nsv` with the provided hook:
  ```{ .sh .no-select }
  nsv tag --hook "./scripts/patch.sh"
  ```
