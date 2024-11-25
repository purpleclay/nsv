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
    hash = "sha256-5B4Vr1D+5EKcDrPcJWssbSLbiX+5CPmzV7qkwkTdph4=";
    leaveDotGit = true;
    postFetch = ''
      cd "$out"
      git rev-parse HEAD > $out/COMMIT
      date -u "+%Y-%m-%dT%H:%M:%SZ" > $out/BUILD_DATE
      find "$out" -name .git -print0 | xargs -0 rm -rf
    '';
  };

  vendorHash = "sha256-acyVQ14yJ2M9YNaq4FNJitZ7J9Kxi45Qj+tqu5GM01c=";

  meta = with lib; {
    homepage = "https://github.com/purpleclay/nsv";
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
