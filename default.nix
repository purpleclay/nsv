{ lib, fetchFromGitHub, buildGoModule }:

let
  version = "0.11.0";
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
    hash = "sha256-/cPNlYgvZuOrD2Kgjmx/8ybLiq+ARXNr7V+TVzqy7JA=";
  };

  vendorHash = "sha256-CR9yD/ksys7rc0jd7Enl14cCDpk0X1Sv84qZGIGw4ak=";

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
