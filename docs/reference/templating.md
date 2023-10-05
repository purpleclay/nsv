---
icon: material/application-cog-outline
status: new
---

# Customization with Templates

Full support for Go [templates](https://pkg.go.dev/text/template) ensures `nsv` is incredibly flexible when generating the next semantic version.

## Next semantic version

The following annotations represent a semantic version broken down into its parts. And are supported by [version customization](../next-version.md#version-template-customization).

| Annotation     | Description                                                             | Example  |
| -------------- | ----------------------------------------------------------------------- | -------- |
| `{{.Prefix}}`  | A monorepo prefix                                                       | `ui/`    |
| `{{.SemVer}}`  | The explicit semantic version. Any leading `v` prefix will be removed   | `0.1.0`  |
| `{{.Version}}` | The version number based on the repositories existing naming convention | `v0.1.0` |

## Tag annotation message :material-new-box:{.new-feature title="Feature added on the 3rd of October 2023"}

The following annotations are available for [customizing](../tag-version.md#using-a-custom-tag-message) the tag annotation message.

| Annotation     | Description                                                      | Example |
| -------------- | ---------------------------------------------------------------- | ------- |
| `{{.Tag}}`     | The next calculated semantic version based on the commit history | `0.2.1` |
| `{{.PrevTag}}` | The previous semantic version                                    | `0.2.0` |
