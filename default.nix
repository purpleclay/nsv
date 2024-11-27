{ lib, fetchFromGitHub, buildGo122Module }:

let
  version = "0.10.1";
in
buildGo122Module {
  pname = "nsv";
  inherit version;
  CGO_ENABLED = 0;

  src = fetchFromGitHub {
    owner = "purpleclay";
    repo = "nsv";
    rev = "v${version}";
    hash = "sha256-HMd6RG0S6ykezV7SC0jLZjxiFExzARTmy8kA3vWuj2g=";
  };

  vendorHash = "sha256-acyVQ14yJ2M9YNaq4FNJitZ7J9Kxi45Qj+tqu5GM01c=";

  meta = with lib; {
    homepage = "https://github.com/purpleclay/nsv";
    changelog = "https://github.com/purpleclay/nsv/releases/tag/v${version}";
    description = "Semantic versioning without any config";
    license = licenses.mit;
    maintainers = with maintainers; [ purpleclay ];
  };

  ldflags = [
    "-s"
    "-w"
    "-X main.version=v${version}"
    "-X main.gitBranch=main"
  ];

  doCheck = false;
}
