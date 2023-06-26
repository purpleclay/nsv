---
icon: material/application-cog-outline
---

# Version Customization with Templates

<span class="rounded-pill">:material-test-tube: experimental</span>

Full support for Go [templates](https://pkg.go.dev/text/template) ensures `nsv` is incredibly flexible when generating the next semantic version.

## Named Fields

The following named fields represent a semantic version broken down into its parts.

| Named Field    | Description                                                             | Example  |
| -------------- | ----------------------------------------------------------------------- | -------- |
| `{{.Prefix}}`  | a monorepo prefix                                                       | `ui/`    |
| `{{.SemVer}}`  | the explicit semantic version. Any leading `v` prefix will be removed   | `0.1.0`  |
| `{{.Version}}` | the version number based on the repositories existing naming convention | `v0.1.0` |
