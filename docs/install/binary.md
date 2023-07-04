---
icon: material/package-variant-closed
---

# Installing the Binary

You can use various package managers to install the `nsv` binary. Take your pick.

## Package Managers

### Homebrew

To use [Homebrew](https://brew.sh/):

```{ .sh .no-select }
brew install purpleclay/tap/nsv
```

### Apt

To install using the [apt](https://ubuntu.com/server/docs/package-management) package manager:

```{ .sh .no-select }
echo 'deb [trusted=yes] https://fury.purpleclay.dev/apt/ /' \
  | sudo tee /etc/apt/sources.list.d/purpleclay.list
sudo apt update
sudo apt install -y nsv
```

You may need to install the `ca-certificates` package if you encounter [trust issues](https://gemfury.com/help/could-not-verify-ssl-certificate/) with regard to the Gemfury certificate:

```{ .sh .no-select }
sudo apt update && sudo apt install -y ca-certificates
```

### Yum

To install using the yum package manager:

```{ .sh .no-select }
echo '[purpleclay]
name=purpleclay
baseurl=https://fury.purpleclay.dev/yum/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/purpleclay.repo
sudo yum install -y nsv
```

### Aur

To install from the [aur](https://archlinux.org/) using [yay](https://github.com/Jguer/yay):

```{ .sh .no-select }
yay -S nsv-bin
```

### Linux Packages

Download and manually install one of the `.deb`, `.rpm` or `.apk` packages from the [Releases](https://github.com/purpleclay/nsv/releases) page.

=== "Apt"

    ```{ .sh .no-select }
    sudo apt install nsv_*.deb
    ```

=== "Yum"

    ```{ .sh .no-select }
    sudo yum localinstall nsv_*.rpm
    ```

=== "Apk"

    ```{ .sh .no-select }
    sudo apk add --no-cache --allow-untrusted nsv_*.apk
    ```

### Go Install

```{ .sh .no-select }
go install github.com/purpleclay/nsv@latest
```

### Bash Script

To install the latest version using a script:

```{ .sh .no-select }
sh -c "$(curl https://raw.githubusercontent.com/purpleclay/nsv/main/scripts/install)"
```

Download a specific version using the `-v` flag. The script uses `sudo` by default but can be disabled through the `--no-sudo` flag. You can also provide a different installation directory from the default `/usr/local/bin` by using the `-d` flag:

```{ .sh .no-select }
sh -c "$(curl https://raw.githubusercontent.com/purpleclay/nsv/main/scripts/install)" \
  -- -v v0.3.0 --no-sudo -d ./bin
```

## Manual download of binary

Head over to the [Releases](https://github.com/purpleclay/nsv/releases) page on GitHub and download any release artefact. Unpack the `nsv` binary and add it to your `PATH`.

## Verifying a binary with cosign

All binaries can be verified using the checksum file and [cosign](https://github.com/sigstore/cosign).

1. Download the checksum file:

   ```sh
   curl -sL https://github.com/purpleclay/nsv/releases/download/v0.3.0/checksums.txt -O
   ```

1. Verify the signature of the file:

   ```sh
   cosign verify-blob \
     --certificate-identity 'https://github.com/purpleclay/nsv/.github/workflows/release.yml@refs/tags/v0.3.0' \
     --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
     --cert 'https://github.com/purpleclay/nsv/releases/download/v0.3.0/checksums.txt.pem' \
     --signature 'https://github.com/purpleclay/nsv/releases/download/v0.3.0/checksums.txt.sig' \
     checksums.txt
   ```

1. Download any release artifact and verify its SHA256 signature matches the entry within the checksum file:

   ```sh
   sha256sum --ignore-missing -c checksums.txt
   ```
