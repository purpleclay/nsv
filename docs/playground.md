---
icon: material/basketball
social:
  cards: false
---

# Explore using the Playground

<span class="rounded-pill">:material-test-tube: experimental</span>

Explore using `nsv` by launching the in-built playground.

## Version templating

Discover how the internal go template is used when generating the next semantic version:

```{ .sh .no-select }
nsv playground ui/v0.1.0 --format '{{.Version}}'
```

```{ .text .no-select .no-copy }
ui/v0.1.0 >> {{.Version}} >> v0.1.0

{{.Prefix}}  >> ui/
{{.SemVer}}  >> 0.1.0
{{.Version}} >> v0.1.0
```

## Command composition

<span class="rounded-pill">:material-bullhorn-variant-outline: coming soon</span>
