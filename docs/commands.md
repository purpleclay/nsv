---
icon: material/apple-keyboard-command
status: new
---

# Control your workflow with Commands

<span class="rounded-pill">:material-test-tube: experimental</span>

Conventional Commits is a great standard, but there are times when it doesn't quite fit with our semantic versioning needs. Commands aim to provide that level of control by including a `nsv: ` footer without your commit message.

Lets start with an example from the Semantic Version 2.0 specification:

```text
Major version zero (0.y.z) is for initial development. Anything MAY
change at any time. The public API SHOULD NOT be considered stable.
```

In a major zero workflow, all changes are deemed unstable. Any conventional commit breaking change should not cause a major increment. As `nsv` is SemVer compliant, a major increment will not happen unless explicitly instructed to do so.

```text
feat: sdfjshdfksjdhfskhfsfksdfhsdjkfhjksdhfkjshfjsdjf

skdfjhsldfhsldkfhskldfjhsdkfhskdjfhskjdfhkjsdfhjsdhfksjdhfkjsdhfksjdhf
sdhfskjdfhskfhskjdhfkjsdfhksjdfhksjdhfjsdhfjkhfdskh

nsv: force~major
```

## Commands

Commands are defined in the footer of a commit message using the case-insensitive `nsv:` prefix. Like conventional commits, it is a simple way of describing intent directly within the commit. And is designed to fit seamlessly within the everyday developer workflow.

### Forcing a semantic increment

`nsv` is designed to scan all commit messages considered part of a release and analyses all defined commands to understand the intent.

- `force~major`
- `force~minor`
- `force~patch`
- `force~ignore`

Ignore disables the force command within the current release, and is used to prevent the need for augmenting the commit history.

### Prerelease support

<span class="rounded-pill">:material-bullhorn-variant-outline: coming soon</span>
