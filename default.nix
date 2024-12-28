{ lib, fetchFromGitHub, buildGoModule }:

let
  version = "0.10.2";
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
    hash = "sha256-tuu7BHWNl6cYlIytBzK1UOWH2kD8C5aA8CDIS1CDYJw=";
  };

  vendorHash = "sha256-ywN9WPEvhfPJfqxXgkq1K7G9fZ9VVIAMqUTN2x0+HMg=";

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
