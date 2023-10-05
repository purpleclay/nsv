---
icon: material/train-car-container
social:
  cards: false
---

# Running with Docker

You can run `nsv` directly from a docker image. Just mount your repository as a volume and set it as the working directory.

=== "DockerHub"

    ```{ .sh .no-select }
    docker run --rm -v $PWD:/work -w /work purpleclay/nsv:v0.3.0
    ```

=== "GHCR"

    ```{ .sh .no-select }
    docker run --rm -v $PWD:/work -w /work ghcr.io/purpleclay/nsv:v0.3.0
    ```

## Verifying with cosign

Docker images can be verified using [cosign](https://github.com/sigstore/cosign).

=== "DockerHub"

    ```{ .sh .no-select }
    cosign verify \
      --certificate-identity 'https://github.com/purpleclay/nsv/.github/workflows/release.yml@refs/tags/v0.3.0' \
      --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
      purpleclay/nsv:v0.3.0
    ```

=== "GHCR"

    ```{ .sh .no-select }
    cosign verify \
      --certificate-identity 'https://github.com/purpleclay/nsv/.github/workflows/release.yml@refs/tags/v0.3.0' \
      --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
      ghcr.io/purpleclay/nsv:v0.3.0
    ```
