---
icon: material/tag-check-outline
status: new
---

# Tag the Next Semantic Version

Let `nsv` tag your repository with the next calculated semantic version:

```sh
nsv tag
```

An annotated tag will be created with the default commit message of `chore: tagged release <version>`.

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

## Signing your tag

If you require your tag to be signed, please ensure your git config is correct before running `nsv`. [gpg-import](https://github.com/purpleclay/gpg-import) is a tool you can easily integrate into your CI workflow and only needs a single environment variable.

## Version template customization

Internally `nsv` utilizes a go template when constructing the next semantic version. Runtime customization of this template is available [here](./next-version.md#version-template-customization).
