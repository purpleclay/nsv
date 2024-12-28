---
icon: material/tools
description: Repair a git repository ready for versioning
---

# Repair a repository ready for versioning

<span class="rounded-pill">:material-test-tube: experimental</span>

Sometimes, a git repository is not in the ideal state for versioning, which can impact how `nsv` behaves.

## Fixing a shallow clone

A git shallow clone limits the amount of history fetched from a repository. It provides performance benefits for day-to-day use but can massively impact `nsv` as it determines the next semantic version. It is <u>especially problematic</u> when working with a monorepo.

Fixing it is a breeze and can be done in one of two ways:

1. You can let `nsv` do it for you:

    === "ENV"

        ```{ .sh .no-select }
        NSV_FIX_SHALLOW=true nsv next
        ```

    === "CLI"

        ```{ .sh .no-select }
        nsv next --fix-shallow
        ```

1. You can fix it through your CI provider of choice. Please refer to their documentation.
