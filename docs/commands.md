---
icon: material/apple-keyboard-command
description: Advanced control of semantic versioning through dedicated commands
---

# Control your workflow with Commands

<span class="rounded-pill">:material-test-tube: experimental</span>

Conventional Commits is an excellent standard, but sometimes it doesn't fit our semantic versioning needs. Let's consider <u>major version zero</u> application development from the Semantic Versioning (`SemVer`) [2.0](https://semver.org/) specification:

```{ .text .no-select .no-copy }
Major version zero (0.y.z) is for initial development. Anything MAY change at
any time. The public API SHOULD NOT be considered stable.
```

All changes within this workflow are not considered stable and could be breaking. Relying purely on the conventional commit standard would cause a major increment for every commit containing a breaking change prefix or footer. As `nsv` is SemVer compliant, a major increment will not happen within this workflow unless explicitly instructed to do so with a command:

```{ .text .no-select .no-copy hl_lines="6" }
feat!: expose new sorting functionality over API

To enable support for server-side sorting, the existing structure of the request
body has been modified to compartmentalize sorting criteria

nsv: force~major
```

## Commands

Commands are defined in the footer of a commit message using the case-insensitive `nsv:` prefix. Like conventional commits, it is a simple way of describing semantic versioning intent and fits seamlessly into any developer's workflow.

### Forcing a semantic increment

The force command allows a developer to take control of the semantic release workflow. Like conventional commits, if multiple commands exist as part of the next release, `nsv` will choose the one that introduces the largest increment.

For simplicity, there is a one-to-one mapping to all semantic increments:

- `force~major`
- `force~minor`
- `force~patch`

If you need to ignore any previous force commands during a release, a break-glass case exists:

- `force~ignore`

### Prerelease support

<span class="rounded-pill">:material-bullhorn-variant-outline: coming soon</span>
