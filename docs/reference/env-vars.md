---
icon: material/earth
status: new
title: Dynamic Configuration
description: Customize behavior with environment variables
---

# Dynamic Configuration with Environment Variables

If you need to customize the behavior of `nsv` you can use the supported environment variables. Environment variables provide a dynamic approach to configuration perfect for integrating `nsv` into your CI workflow.

## Variables

| Variable Name        | Description                                                                                                   |
| -------------------- | ------------------------------------------------------------------------------------------------------------- |
| `LOG_LEVEL`          | the level of logging when printing to stderr <br/>(`debug`, `info`, `warn`, `error`, `fatal`)                 |
| `NO_COLOR`           | switch to using an ASCII color profile within the terminal                                                    |
| `NO_LOG`             | disable all log output                                                                                        |
| `NSV_FORMAT`         | set a go template for formatting the provided tag                                                             |
| `NSV_MAJOR_PREFIXES` | a comma separated list of conventional commit prefixes for triggering <br/>a major semantic version increment |
| `NSV_MINOR_PREFIXES` | a comma separated list of conventional commit prefixes for triggering <br/>a minor semantic version increment |
| `NSV_PATCH_PREFIXES` | a comma separated list of conventional commit prefixes for triggering <br/>a patch semantic version increment |
| `NSV_PRETTY`         | pretty-print the output of the next semantic version in a given format                                        |
| `NSV_SHOW`           | show how the next semantic version was generated                                                              |
| `NSV_TAG_MESSAGE`    | a custom message for the tag, overrides the default:<br/>`chore: tagged release <version>`                    |
