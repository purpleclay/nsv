---
icon: material/console
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
| NSV_COMMIT_MESSAGE | a custom message when committing file changes,         |
|                    | supports go text templates. The default is: "chore:    |
|                    | patched files for release {{.Tag}}                     |
|                    | {{.SkipPipelineTag}}"                                  |
| NSV_DRY_RUN        | no changes will be made to the repository              |
| NSV_FIX_SHALLOW    | fix a shallow clone of a repository if detected        |
| NSV_FORMAT         | provide a go template for changing the default version |
|                    | format                                                 |
| NSV_HOOK           | a user-defined hook that will be executed before the   |
|                    | repository is tagged with the next semantic version    |
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

Hook Environment Variables:

| Name                  | Description                                         |
|-----------------------|-----------------------------------------------------|
| NSV_NEXT_TAG          | the next calculated semantic version                |
| NSV_PREV_TAG          | the last semantic version as identified within the  |
|                       | tag                                                 |
|                       | history of the current repository                   |
| NSV_WORKING_DIRECTORY | the working directory (or path) relative to the     |
|                       | root of the current repository. It will be empty if |
|                       | not a monorepo                                      |
```

## Usage

```{ .text .no-select .no-copy }
nsv tag [<path>...] [flags]
```

## Flags

```{ .text .no-select .no-copy }
-M, --commit-message string    a custom message when committing file changes,
                               supports go text templates (default "chore:
                               tagged release {{.Tag}} {{.SkipPipelineTag}}")
    --dry-run                  no changes will be made to the repository
    --fix-shallow              fix a shallow clone of a repository if detected
-f, --format string            provide a go template for changing the default
                               version format
-h, --help                     help for tag
    --hook string              a user-defined hook that will be executed before
                               the repository is tagged with the next semantic
                               version
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
-A, --tag-message string       a custom message for the annotated tag, supports
                               go text templates (default "chore: tagged release
                               {{.Tag}}")
```

## Global Flags

```{ .text .no-select .no-copy }
--log-level string   the level of logging when printing to stderr
                     (default "info")
--no-color           switch to using an ASCII color profile within the terminal
--no-log             disable all log output
```
