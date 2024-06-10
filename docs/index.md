---
title: Next Semantic Version
description: Convention-based semantic versioning without a config file
---

# NSV

`nsv` (Next Semantic Version) is a convention-based semantic versioning tool that leans on the power of conventional commits to make versioning your software a breeze!

## See it in action

<div>
    <video controls>
        <source src="./static/nsv.webm" type="video/webm">
        <source src="./static/nsv.mp4" type="video/mp4">
    </video>
</div>

## Why another versioning tool

There are many semantic versioning tools already out there! But they typically require some configuration or custom scripting in your CI system to make them work. No one likes managing config; it is error-prone, and the slightest tweak ultimately triggers a cascade of change across your projects.

`nsv` makes intelligent semantic versioning decisions about your project without needing a config file. It is convention-based and adapts to your semantic workflow by analyzing your commit messages.

The power is at your fingertips.

### Features

- First-class support for <u>semantic versioning</u>. Conventional commits help give it a nudge in the right direction.
- Configure your Conventional Commit prefixes for complete control.
- If you batch your commits per release or prefer a continuous delivery approach, it has you covered.
- <u>Context-aware</u>, it automatically switches to a monorepo workflow.
- Extend the power of your commits through <u>commands</u> to dynamically change your semantic release workflow.
- Explore how to use it within the purpose-built playground. <span class="rounded-pill">:material-test-tube: experimental</span>
- Get up and running in seconds within GitHub :material-github: and GitLab :material-gitlab: with the available <u>[action](https://github.com/purpleclay/nsv-action)</u> or <u>[template](https://gitlab.com/purpleclay/nsv)</u>.
- First class <u>[Dagger](https://dagger.io/)</u> :material-sword: support. <span class="rounded-pill">:material-test-tube: experimental</span>
