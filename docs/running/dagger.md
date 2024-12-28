---
icon: material/sword
social:
  cards: false
---

# Running with Dagger

[Dagger](https://dagger.io/) provides a way to define tasks that can run on your favorite CI/CD platform or locally. Powered by Docker, it is incredibly easy to get up and running.

## Installing Dagger

Read the official documentation for complete [instructions](https://docs.dagger.io/install).

=== "macOS"

    ```{ .sh .no-select }
    brew install dagger/tap/dagger
    ```

=== "Linux"

    ```{ .sh .no-select }
    curl -L https://dl.dagger.io/dagger/install.sh | BIN_DIR=$HOME/.local/bin sh
    ```

=== "Windows"

    ```{ .sh .no-select }
    Invoke-WebRequest -UseBasicParsing -Uri https://dl.dagger.io/dagger/install.ps1 | Invoke-Expression
    ```

=== "Nix Flake"

    ```{ .nix .no-select }
    {
      inputs = {
        nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
        flake-utils.url = "github:numtide/flake-utils";

        dagger = {
          url = "github:dagger/nix";
          inputs = {
            nixpkgs.follows = "nixpkgs";
          };
        };
      };

      outputs = { self, nixpkgs, flake-utils, dagger }:
        flake-utils.lib.eachDefaultSystem (system:
          let
            pkgs = nixpkgs.legacyPackages.${system};
          in
          with pkgs;
          {
            devShells.default = mkShell {
              buildInputs = [
                dagger.packages.${system}.dagger
              ];
            };
          }
        );
    }
    ```

## Next semantic version

Calculates the next semantic version and prints it to stdout. For full configuration [options](../next-version.md).

```{ .sh .no-select }
dagger call -m github.com/purpleclay/daggerverse/nsv@v0.10.1 --src . next
```

## Tagging the next version

Tags a repository with the next semantic version. For full configuration [options](../tag-version.md).

```{ .sh .no-select }
dagger call -m github.com/purpleclay/daggerverse/nsv@v0.10.1 --src . tag
```

## Patching files with the next version

Patch files within a repository using the next semantic version. For full configuration [options](../patch-files.md).

```{ .sh .no-select }
dagger call -m github.com/purpleclay/daggerverse/nsv@v0.10.1 --src . patch \
  --hook "./scripts/patch.sh"
```
