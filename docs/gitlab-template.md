---
icon: material/gitlab
---

# Using the GitLab template

To get up and running within a GitLab pipeline, include the publicly available `nsv` GitLab [template](https://gitlab.com/purpleclay/nsv). You can find details on setting `environment variables` in the documentation.

## Tagging a repository

```{.yaml linenums="1"}
include:
  - https://gitlab.com/purpleclay/nsv/-/raw/main/nsv.gitlab-ci.yml
```
