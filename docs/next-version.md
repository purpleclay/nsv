---
icon: material/tag-arrow-up-outline
status: new
---

# Next Semantic Version

NSV core principles of being context-aware and convention-based without a config file let you achieve most of your semantic versioning needs by simply running:

```sh
nsv next
```

By scanning all commit messages within the latest release, `nsv` understands the author's intent and prints the next semantic version to stdout.

If you want to see what is happening under the hood, ask it:

```sh
nsv next --show
```

```text
EXAMPLE OF OUTPUT
```

If you need to customize its behavior, environment variables, CLI flags, or [commands](./commands.md) can be used.

## Monorepos as first-class citizens

Monorepo support is not an afterthought. By being context-aware, `nsv` can detect if it runs outside the repository root and calculates the next semantic version based on its location.

```text
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

```sh
{{.Prefix}}{{.Version}}
```

Runtime customization of this template is available. For example, you can enforce explicit semantic version usage:

=== "ENV"

    ```sh
    NSV_FORMAT="{{.SemVer}}" nsv next
    ```

=== "CLI"

    ```sh
    nsv next --format "{{.SemVer}}"
    ```

Head over to the [playground](./playground.md) to discover more.

## Prerelease support

<span class="rounded-pill">:material-bullhorn-variant-outline: coming soon</span>
