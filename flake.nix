{
  description = "Semantic versioning without any config";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    # Lock to version: 0.14.0
    dagger = {
      url = "github:dagger/nix?rev=9852fdddcdcb52841275ffb6a39fa1524d538d5a";
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
            git
            go
            go-task
            goreleaser
            vhs
          ];
        };

        apps.default = {
          type = "app";
          program = "${pkgs.callPackage ./. {}}/bin/nsv";
        };

        packages.default = pkgs.callPackage ./. {};
      }
    );
}
