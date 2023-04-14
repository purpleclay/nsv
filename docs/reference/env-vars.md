---
icon: material/earth
status: new
---

# Dynamic Configuration with Environment Variables

<span class="rounded-pill">:material-test-tube: experimental</span>

If you need to customize the behavior of `nsv` you can use the supported environment variables. Environment variables provide a dynamic approach to configuration perfect for integrating `nsv` into your CI workflow.

## Variables

| Variable Name | Description                                       |
| ------------- | ------------------------------------------------- |
| `NSV_FORMAT`  | set a go template for formatting the provided tag |
