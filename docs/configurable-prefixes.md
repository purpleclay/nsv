---
icon: material/apple-keyboard-command
description: Advanced control of semantic versioning through dedicated commands
---

# Configure your Conventional prefixes

[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) offers a fantastic approach to controlling semantic versioning through commit prefixes. An industry standard convention devised by the Angular team is most commonly used to date.

- `BREAKING CHANGE`: any commit with this footer triggers a major update.
- `!`: any conventional prefix with this suffix (e.g. `refactor!`) triggers a major update.
- `feat`: triggers a minor update.
- `fix`: triggers a patch update.

## Defining your own rules

- `breaking`: triggers a major update `1.0.0` ~> `2.0.0`.
- `feat`, `deps`: triggers a minor update `0.1.0` ~> `0.2.0`.
- `fix`, `docs`, `styles`: triggers a patch update `0.3.1` ~> `0.3.2`.

=== "ENV"

    ```{ .sh .no-select }
    NSV_MAJOR_PREFIXES=breaking \
      NSV_MINOR_PREFIXES=feat,deps \
      NSV_PATCH_PREFIXES=fix,docs,styles \
      nsv next
    ```

=== "CLI"

    ```{ .sh .no-select }
    nsv next --major-prefixes breaking \
      --minor-prefixes feat,deps \
      --patch-prefixes fix,docs,styles
    ```

Don't worry—when defining your custom prefixes, both the `BREAKING CHANGE` footer and the `!` suffix are automatically supported.

## How prefix matching works

`nsv` matches a prefix in one of two ways:

- `breaking` is a <u>wildcard</u> prefix capable of matching against an optional scope.
- `breaking(api)` is an <u>exact</u> match only.
