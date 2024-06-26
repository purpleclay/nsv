---
icon: material/apple-keyboard-command
description: Advanced control of semantic versioning through dedicated commands
---

# Control your workflow with Commands

Conventional Commits is an excellent standard, but sometimes it doesn't fit our semantic versioning needs. Let's consider <u>major version zero</u> application development from the Semantic Versioning (`SemVer`) [2.0](https://semver.org/) specification:

```{ .text .no-select .no-copy }
Major version zero (0.y.z) is for initial development. Anything MAY change at
any time. The public API SHOULD NOT be considered stable.
```

Within this workflow, the largest semantic increase would be a `minor` increment. Relying purely on the conventional commit standard will not trigger a major increment. As `nsv` is SemVer compliant, it too will adhere to this, unless explicitly instructed to do so with a command:

```{ .text .no-select .no-copy hl_lines="6" }
feat!: expose new sorting functionality over API

To enable support for server-side sorting, the existing structure of the request
body has been modified to compartmentalize sorting criteria

nsv: force~major
```

## Commands

Commands are defined in the footer of a commit message using the case-insensitive `nsv:` prefix. Like conventional commits, it is a simple way of describing semantic versioning intent and is designed to fit seamlessly into any developer's workflow. When searching for commands, `nsv` <u>**stops at the first one it finds**</u>. This is an important distinction between using commands and conventional commits to control SemVer.

Multiple commands can be grouped together to achieve your desired outcome:

```{ .go .annotate .no-select .no-copy }
nsv: force~major,pre // (1)!
```

1. If the existing SemVer was `< 1.0.0`, these commands combined would generate a SemVer of `1.0.0-beta.1`. Read on to discover how these commands work.

### Forcing a semantic increment

The `force` command allows a developer to take complete control over the semantic release workflow, ignoring any existing conventional commits. It consists of two parts, a mandatory label `force~`, followed by the desired semantic increment:

```{ .text .no-select .no-copy }
nsv: force~major
```

For simplicity, there is a one-to-one mapping to all SemVer increments:

- `force~major`
- `force~minor`
- `force~patch`

If you need to ignore any previous force commands, a break-glass command exists:

- `force~ignore`

### Prerelease support :material-new-box:{.new-feature title="Feature added on the 21st of February 2024"}

<span class="rounded-pill">:material-test-tube: experimental</span>

The `pre` command allows a developer to initiate a semantic prerelease workflow, which isn't possible through conventional commits. It consists of two parts, a mandatory label `pre`, followed by an optional prerelease version:

```{ .text .no-select .no-copy }
nsv: pre~alpha
```

There is a one-to-one mapping to common prerelease labels:

- `pre~alpha`
- `pre~beta`
- `pre~rc`
- `pre` on its own is equivalent to `pre~beta`

A prerelease version generated by the `pre` command follows the SemVer convention of:

```{ .text .no-select .no-copy }
0.1.0-beta.1
```

The `.1` part of the version is automatically incremented by `nsv` for each subsequent SemVer prerelease. It is reset when transitioning between prerelease labels.
