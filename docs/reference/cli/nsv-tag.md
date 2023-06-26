---
icon: material/console
status: new
---

# nsv tag

```{ .text .no-select .no-copy }
Tag the repository with the next semantic version based on the conventional
commit history of your repository.

Environment Variables:

| Name            | Description                                              |
|-----------------|----------------------------------------------------------|
| NSV_FORMAT      | provide a go template for changing the default version   |
|                 | format                                                   |
| NSV_SHOW        | show how the next semantic version was generated         |
| NSV_TAG_MESSAGE | a custom message for the tag, overrides the default      |
|                 | message of: chore: tagged <version> by nsv               |
```

## Usage

```{ .text .no-select .no-copy }
nsv tag [flags]
```

## Flags

```{ .text .no-select .no-copy }
-f, --format string    provide a go template for changing the default version
                       format
-h, --help             help for tag
-m, --message string   a custom message for the tag, overrides the default
                       message of: chore: <version> tagged by nsv
-s, --show             show how the next semantic version was generated
```
