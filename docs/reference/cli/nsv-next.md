---
icon: material/console
---

# nsv next

```{ .text .no-select .no-copy }
Generate the next semantic version based on the conventional commit history
of your repository.

Environment Variables:

| Name       | Description                                                   |
|------------|---------------------------------------------------------------|
| NSV_FORMAT | provide a go template for changing the default version format |
| NSV_SHOW   | show how the next semantic version was generated              |
```

## Usage

```{ .text .no-select .no-copy }
nsv next [flags]
```

## Flags

```{ .text .no-select .no-copy }
-f, --format string   provide a go template for changing the default version
                      format
-h, --help            help for next
-s, --show            show how the next semantic version was generated
```
