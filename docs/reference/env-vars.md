---
icon: material/earth
status: new
---

# Dynamic Configuration with Environment Variables

If you need to customize the behavior of `nsv` you can use the supported environment variables. Environment variables provide a dynamic approach to configuration perfect for integrating `nsv` into your CI workflow.

## Variables

| Variable Name     | Description                                                                            |
| ----------------- | -------------------------------------------------------------------------------------- |
| `NO_COLOR`        | switch to using an ASCII color profile within the terminal                             |
| `NSV_FORMAT`      | set a go template for formatting the provided tag                                      |
| `NSV_PRETTY`      | pretty-print the output of the next semantic version in a given format                 |
| `NSV_SHOW`        | show how the next semantic version was generated                                       |
| `NSV_TAG_MESSAGE` | a custom message for the tag, overrides the default: `chore: tagged release <version>` |
