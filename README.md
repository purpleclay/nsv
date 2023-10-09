# NSV

NSV (Next Semantic Version) is a convention-based semantic versioning tool that leans on the power of conventional commits to make versioning your software a breeze!

## See it in action

https://github.com/purpleclay/nsv/assets/106762954/4d3362ac-6b1c-42b1-a13a-0edffb2ab6bd

## Why another versioning tool

There are many semantic versioning tools already out there! But they typically require some configuration or custom scripting in your CI system to make them work. No one likes managing config; it is error-prone, and the slightest tweak ultimately triggers a cascade of change across your projects.

Step in NSV. Designed to make intelligent semantic versioning decisions about your project without needing a config file. Entirely convention-based, you can adapt your workflow from within your commit message.

The power is at your fingertips.

### Features

- First-class support for semantic versioning. Conventional commits help give it a nudge in the right direction.
- If you batch your commits per release or prefer a continuous delivery approach, it has you covered.
- Context-aware, it automatically switches to a monorepo workflow.
- Extend the power of your commits through commands to dynamically change your semantic release workflow.
- Explore how to use it within the purpose-built playground.
- Get up and running in seconds within GitHub and GitLab with the available [action](https://github.com/purpleclay/nsv-action) or [template](https://gitlab.com/purpleclay/nsv).

## Documentation

Check out the latest [documentation](https://docs.purpleclay.dev/nsv/)

## Badges

[![Build status](https://img.shields.io/github/actions/workflow/status/purpleclay/nsv/ci.yml?style=flat-square&logo=go)](https://github.com/purpleclay/nsv/actions?workflow=ci)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/purpleclay/nsv.svg?style=flat-square)](go.mod)
[![DeepSource](https://app.deepsource.com/gh/purpleclay/nsv.svg/?label=active+issues&token=1JHae42sEcDRhZdEe-tfiPEv)](https://app.deepsource.com/gh/purpleclay/nsv/?ref=repository-badge)
