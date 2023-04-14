---
icon: material/file-code-outline
---

# Compiling from Source

Download both [Go 1.20+](https://go.dev/doc/install) and [go-task](https://taskfile.dev/#/installation). Then clone the code from GitHub:

```sh
git clone https://github.com/purpleclay/nsv.git
cd nsv
```

Build:

```sh
task
```

And check that everything works:

```sh
./nsv version
```
