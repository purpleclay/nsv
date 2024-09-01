---
icon: material/application-cog-outline
title: Flexibility through Templating
description: Customize behavior with go templates
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

## Tag annotation message

The following annotations are available for [customizing](../tag-version.md#using-a-custom-tag-message) the tag annotation message.

| Annotation     | Description                                                      | Example |
| -------------- | ---------------------------------------------------------------- | ------- |
| `{{.Tag}}`     | The next calculated semantic version based on the commit history | `0.2.1` |
| `{{.PrevTag}}` | The previous semantic version                                    | `0.2.0` |

## Commit message

The following annotations are available for [customizing](../tag-version.md#using-a-custom-commit-message) the commit message.

| Annotation             | Description                                                      | Example     |
| ---------------------- | ---------------------------------------------------------------- | ----------- |
| `{{.Tag}}`             | The next calculated semantic version based on the commit history | `0.2.1`     |
| `{{.PrevTag}}`         | The previous semantic version                                    | `0.2.0`     |
| `{{.SkipPipelineTag}}` | A CI provider tag for skipping a pipeline build                  | `[skip ci]` |
