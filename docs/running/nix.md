---
icon: material/nix
social:
  cards: false
---

# Running with Nix

You can run `nsv` directly using [Nix](https://zero-to-nix.com/concepts/flakes/) without directly installing its dependencies.

## Installing Nix

Determinate Systems provides a quick and easy [installer](https://github.com/DeterminateSystems/nix-installer).

```{ .sh .no-select }
curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | \
  sh -s -- install
```

## Nix Run

With a single command, Nix will build and run `nsv` within your environment:

```sh
nix run github:purpleclay/nsv
```

Passing command line arguments is also very easy:

```sh
nix run github:purpleclay/nsv -- next
```
