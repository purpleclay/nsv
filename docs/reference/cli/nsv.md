---
icon: material/console
---

# nsv

```{ .text .no-select .no-copy }
NSV (Next Semantic Version) is a convention-based semantic versioning tool that
leans on the power of conventional commits to make versioning your software a
breeze!.

## Why another versioning tool

There are many semantic versioning tools already out there! But they typically
require some configuration or custom scripting in your CI system to make them
work. No one likes managing config; it is error-prone, and the slightest tweak
ultimately triggers a cascade of change across your projects.

Step in NSV. Designed to make intelligent semantic versioning decisions about
your project without needing a config file. Entirely convention-based, you can
adapt your workflow from within your commit message.

The power is at your fingertips.

Global Environment Variables:

| Name     | Description                                                   |
|----------|---------------------------------------------------------------|
| NO_COLOR | switch to using an ASCII color profile within the terminal    |
```

## Usage

```{ .text .no-select .no-copy }
nsv [command]
```

## Commands

```{ .text .no-select .no-copy }
completion  Generate the autocompletion script for the specified shell
help        Help about any command
next        Generate the next semantic version
playground  A playground for discovering go template support
tag         Tag the repository with the next semantic version
version     Print build time version information
```

## Flags

```{ .text .no-select .no-copy }
-h, --help       help for nsv
    --no-color   switch to using an ASCII color profile within the terminal
```
