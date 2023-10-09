---
icon: material/tag-arrow-up-outline
status: new
description: Determine the next semantic version of your repository based on its history
---

# Next Semantic Version

`nsv` core principles of being <u>context-aware</u> and <u>convention-based</u> will let you achieve almost all of your semantic versioning needs when running:

```{ .sh .no-select }
nsv next
```

By scanning all commit messages within the latest release, `nsv` understands the author's intent and prints the next semantic version to stdout.

If you want to see what is happening under the hood:

=== "ENV"

    ```{ .sh .no-select }
    NSV_SHOW=true nsv next
    ```

=== "CLI"

    ```{ .sh .no-select }
    nsv next --show
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

If you need to customize its behavior, [environment variables](./reference/env-vars.md), CLI [flags](./reference/cli/nsv-next.md), or [commands](./commands.md) can be used.

## Configurable paths for monorepo support :material-new-box:{.new-feature title="Feature added on the 9th of October 2023"}

Monorepo support is important to the [design](./monorepos.md) of `nsv`. By adding support for context paths, multiple semantic versions can be resolved throughout a repository in a single operation.

```{ .sh .no-select }
nsv next src/ui src/search
```

Any version change will be printed to stdout as a comma separated list in context path order:

```{ .text .no-select .no-copy }
ui/0.2.1,search/0.3.0
```

## Version template customization

Internally `nsv` utilizes a go template when constructing the next semantic version:

```{ .sh .no-select }
{{.Prefix}}{{.Version}}
```

Runtime customization of this template is available. For example, you can enforce explicit semantic version usage:

=== "ENV"

    ```{ .sh .no-select }
    NSV_FORMAT="{{.SemVer}}" nsv next
    ```

=== "CLI"

    ```{ .sh .no-select }
    nsv next --format "{{.SemVer}}"
    ```

Head over to the [playground](./playground.md) to discover more.
