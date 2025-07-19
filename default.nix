{ lib, fetchFromGitHub, buildGoModule }:

let
  version = "0.12.1";
in
buildGoModule {
  pname = "nsv";
  inherit version;
  env.CGO_ENABLED = 0;

  src = fetchFromGitHub {
    owner = "purpleclay";
    repo = "nsv";
    rev = "v${version}";
    leaveDotGit = true;
    postFetch = ''
      cd "$out"
      git rev-parse HEAD > $out/COMMIT
      git log -1 --format=%cI > $out/BUILD_DATE
      find "$out" -name .git -print0 | xargs -0 rm -rf
    '';
    hash = "sha256-qhU8ONGprrgBL2gYQj/yN7iRf0v2tod73ThysmU4AfM=";
  };

  vendorHash = "sha256-x5R2eAGEfVZWo+eh7KwMd0Q5aKx95qgVvrOrT6suh8g=";

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

  preBuild = ''
    ldflags+=" -X main.gitCommit=$(cat COMMIT) -X main.buildDate=$(cat BUILD_DATE)"
  '';

  doCheck = false;
}
