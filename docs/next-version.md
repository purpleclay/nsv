---
icon: material/tag-arrow-up-outline
---

# Next Semantic Version

`nsv` core principles of being <u>context-aware</u> and <u>convention-based</u> will let you achieve almost all of your semantic versioning needs when running:

```{ .sh .no-select }
nsv next
```

By scanning all commit messages within the latest release, `nsv` understands the author's intent and prints the next semantic version to stdout.

If you want to see what is happening under the hood you can use the `--show` flag:

```{ .sh .no-select }
nsv next --show
```

```{ .text .no-select .no-copy }
0.1.0

 HEAD  ...  0.1.0
────────────────────────────────────────────────────
c6bfdda fix: fix to the store
a0a1e2b feat: new exciting search feature << matched
83def28 ci: configure workflows
b8a7daf chore: scaffold project
6c05c93 initialize repository
────────────────────────────────────────────────────
```

If you need to customize its behavior further, [environment variables](./reference/env-vars.md), CLI [flags](./reference/cli/nsv-next.md), or [commands](./commands.md) can be used.

## Monorepos as first-class citizens

Monorepo support is not an afterthought. By being context-aware, `nsv` can detect if it runs outside the repository root and calculates the next semantic version based on its location.

```{ .text .no-select .no-copy }
awesome-app
  ui           << ui/v0.2.1
    ...
  backend      << backend/v0.3.0
    ...
  store        << store/v0.2.3
    ...
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

## Prerelease support

<span class="rounded-pill">:material-bullhorn-variant-outline: coming soon</span>
