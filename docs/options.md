---
icon: material/cog-outline
social:
  cards: false
status: new
---

# Different ways of customizing NSV

`nsv` is designed to be config free and require minimal to no runtime options to release your software. If an option must be set, you have one of <u>three ways</u> to set them, ordered from highest to lowest precedence.

## CLI flags

A CLI flag is the most explicit way to set an option and will always override the other two methods.

```{ .sh .no-select }
nsv next --format "{{.SemVer}}" --fix-shallow
```

## Environment variables

Environment variables provide great flexibility and support dynamic setting of options, for example, within a CI pipeline.

```{ .sh .no-select }
NSV_FORMAT="{{.SemVer}}" NSV_FIX_SHALLOW="true" nsv next
```

## DotEnv file

A `.env` file located in the root of a project is automatically loaded at runtime and injected into the `nsv` context.

```{ .sh .no-select }
NSV_FORMAT="{{.SemVer}}"
NSV_FIX_SHALLOW="true"
```
