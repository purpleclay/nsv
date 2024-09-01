---
icon: material/console
social:
  cards: false
---

# nsv next

```{ .text .no-select .no-copy }
Generate the next semantic version based on the conventional commit history
of your repository.

Environment Variables:

| Name               | Description                                            |
|--------------------|--------------------------------------------------------|
| LOG_LEVEL          | the level of logging when printing to stderr           |
|                    | (default: info)                                        |
| NO_COLOR           | switch to using an ASCII color profile within the      |
|                    | terminal                                               |
| NO_LOG             | disable all log output                                 |
| NSV_FORMAT         | provide a go template for changing the default version |
|                    | format                                                 |
| NSV_MAJOR_PREFIXES | a comma separated list of conventional commit prefixes |
|                    | or triggering a major semantic version increment       |
| NSV_MINOR_PREFIXES | a comma separated list of conventional commit prefixes |
|                    | for triggering a minor semantic version increment      |
| NSV_PATCH_PREFIXES | a comma separated list of conventional commit prefixes |
|                    | for triggering a patch semantic version increment      |
| NSV_PRETTY         | pretty-print the output of the next semantic version   |
|                    | in a given format. The format can be one of either     |
|                    | full or compact. Must be used in conjunction with      |
|                    | NSV_SHOW (default: full)                               |
| NSV_SHOW           | show how the next semantic version was generated       |
```

## Usage

```{ .text .no-select .no-copy }
nsv next [<path>...] [flags]
```

## Flags

```{ .text .no-select .no-copy }
-f, --format string            provide a go template for changing the default
                               version format
-h, --help                     help for next
    --major-prefixes strings   a comma separated list of conventional commit
                               prefixes for triggering a major semantic version
                               increment
    --minor-prefixes strings   a comma separated list of conventional commit
                               prefixes for triggering a minor semantic version
                               increment
    --patch-prefixes strings   a comma separated list of conventional commit
                               prefixes for triggering a patch semantic version
                               increment
-p, --pretty string            pretty-print the output of the next semantic
                               version in a given format. The format can be one
                               of either full or compact. Must be used in
                               conjunction with --show (default "full")
-s, --show                     show how the next semantic version was generated
```

## Global Flags

```{ .text .no-select .no-copy }
--log-level string   the level of logging when printing to stderr
                     (default "info")
--no-color           switch to using an ASCII color profile within the terminal
--no-log             disable all log output
```
