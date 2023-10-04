---
icon: material/console
status: new
---

# nsv tag

```{ .text .no-select .no-copy }
Tag the repository with the next semantic version based on the conventional
commit history of your repository.

Environment Variables:

| Name            | Description                                               |
|-----------------|-----------------------------------------------------------|
| NO_COLOR        | switch to using an ASCII color profile within the         |
|                 | terminal                                                  |
| NSV_FORMAT      | provide a go template for changing the default version    |
|                 | format                                                    |
| NSV_PRETTY      | pretty-print the output of the next semantic version in   |
|                 | a given format. The format can be one of either full or   |
|                 | compact. full is the default. Must be used in conjunction |
|                 | with NSV_SHOW                                             |
| NSV_SHOW        | show how the next semantic version was generated          |
| NSV_TAG_MESSAGE | a custom message for the tag, supports go text templates. |
|                 | The default is: "chore: tagged release {{.Tag}}"          |
```

## Usage

```{ .text .no-select .no-copy }
nsv tag [<path>...] [flags]
```

## Flags

```{ .text .no-select .no-copy }
-f, --format string    provide a go template for changing the default version
                       format
-h, --help             help for tag
-m, --message string   a custom message for the tag, supports go text templates
                       (default "chore: tagged release {{.Tag}}")
-p, --pretty string    pretty-print the output of the next semantic version in
                       a given format. The format can be one of either full or
                       compact. Must be used in conjunction with --show
                       (default "full")
-s, --show             show how the next semantic version was generated
```

## Global Flags

```{ .text .no-select .no-copy }
--no-color   switch to using an ASCII color profile within the terminal
```
