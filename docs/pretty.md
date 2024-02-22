---
icon: material/panorama-variant
description: Pretty-printing formats for customizing output
---

# Pretty-Printing Output

Support for pretty-printing allows the output of `nsv` to be customized. A set of in-built options are supported, with more to follow. The default option is `full` but can be changed:

=== "ENV"

    ```{ .sh .no-select }
    NSV_PRETTY=compact nsv next
    ```

=== "CLI"

    ```{ .sh .no-select }
    nsv next --pretty compact
    ```

!!! tip "A colorless output"

    Stripping back to basics is easy, set the following environment variable to switch to an ASCII color profile, `NO_COLOR=1`.

## Full

A tabular format displaying each semantic version change with its associated history, clearly highlighting the triggering commit.

## Compact

A tabular format with a condensed history, providing a focused overview of any triggering commit.

```{ .text .no-select .no-copy }
┌───────────────┬──────────────────────────────────────────────────┐
│  0.2.0        │ ✓ 2020953                                        │
│  ↑↑           │   >>feat<<: a new exciting feature               │
│  0.1.0        │                                                  │
└───────────────┴──────────────────────────────────────────────────┘
```
