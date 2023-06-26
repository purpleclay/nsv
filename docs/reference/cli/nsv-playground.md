---
icon: material/console
---

# nsv playground

```{ .text .no-select .no-copy }
A playground for discovering go template support.

Discover ways of formatting your repository tag using the in-built
go template annotations.

Environment Variables:

| Name       | Description                                       |
|------------|---------------------------------------------------|
| NSV_FORMAT | set a go template for formatting the provided tag |
```

## Usage

```{ .text .no-select .no-copy }
nsv playground <tag> [flags]
```

## Flags

```{ .text .no-select .no-copy }
-f, --format string   provide a go template for changing the default version
                      format
-h, --help            help for playground
```
