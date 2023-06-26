---
icon: material/earth
status: new
---

# Dynamic Configuration with Environment Variables

<span class="rounded-pill">:material-test-tube: experimental</span>

If you need to customize the behavior of `nsv` you can use the supported environment variables. Environment variables provide a dynamic approach to configuration perfect for integrating `nsv` into your CI workflow.

## Variables

| Variable Name     | Description                                                                            |
| ----------------- | -------------------------------------------------------------------------------------- |
| `NSV_FORMAT`      | set a go template for formatting the provided tag                                      |
| `NSV_SHOW`        | show how the next semantic version was generated                                       |
| `NSV_TAG_MESSAGE` | a custom message for the tag, overrides the default: `chore: tagged release <version>` |
