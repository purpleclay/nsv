---
icon: material/dots-hexagon
status: new
---

# Monorepos as first-class citizens

`nsv` has been designed so that monorepo support is <u>not an afterthought</u>. By being context-aware, `nsv` can detect if it runs outside the repository root and calculates the next semantic version based on its location.

```{ .text .no-select .no-copy }
awesome-app
  ui           << ui/v0.2.1
    ...
  backend      << backend/v0.3.0
    ...
  store        << store/v0.2.3
    ...
```

1. context aware - understands its location outside of the repository root
1. user specified paths can be provided

## Understands its running context

example of running nsv within a different directory

```sh
nsv next
```

## A

example of running nsv with a path

```sh
nsv tag
```
