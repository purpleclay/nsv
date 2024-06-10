---
icon: material/console
status: new
social:
  cards: false
---

# nsv tag

```{ .text .no-select .no-copy }
Tag the repository with the next semantic version based on the conventional
commit history of your repository.

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
| NSV_TAG_MESSAGE    | a custom message for the tag, supports go text         |
|                    | templates. The default is: "chore: tagged release      |
|                    | {{.Tag}}"                                              |
```

## Usage

```{ .text .no-select .no-copy }
nsv tag [<path>...] [flags]
```

## Flags

```{ .text .no-select .no-copy }
-f, --format string            provide a go template for changing the default
                               version format
-h, --help                     help for tag
    --major-prefixes strings   a comma separated list of conventional commit
                               prefixes for triggering a major semantic version
                               increment
-m, --message string           a custom message for the tag, supports go text
                               templates (default "chore: tagged release
                               {{.Tag}}")
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
