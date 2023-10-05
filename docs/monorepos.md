---
icon: material/dots-hexagon
status: new
---

# Monorepos as first-class citizens

`nsv` has been designed so that monorepo support is <u>not an afterthought</u>. Monorepo detection is built-in, removing the need for additional configuration.

## Understands its running context

By being context-aware, `nsv` can detect if it runs within a repository subdirectory, changing how it inspects the commit history. The next semantic version will include the component prefix, a standard monorepo practice[^1].

```{ .sh .no-select .no-copy }
cd src/ui
```

```{ .sh .no-select .no-copy }
$ nsv next

ui/0.2.0
```

Context paths as command line arguments remove the need to change directories. `nsv` can version multiple monorepo components in a single pass.

```{ .sh .no-select .no-copy }
$ nsv next src/ui src/search src/database

ui/0.3.0,search/0.2.1,database/0.3.0
```

[^1]: Full [customization](./next-version.md#version-template-customization) is supported through Go templating if you want to change this behavior.
