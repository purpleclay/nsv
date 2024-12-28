# syntax=docker/dockerfile:1
FROM alpine:3.21

# renovate: datasource=github-tags depName=purpleclay/gpg-import versioning=semver
ENV GPG_IMPORT_VERSION=0.4.0

RUN apk add --no-cache git git-lfs gnupg tini curl
RUN sh -c "$(curl https://raw.githubusercontent.com/purpleclay/gpg-import/main/scripts/install)" -- -v ${GPG_IMPORT_VERSION}

ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
CMD ["--help"]

COPY --chmod=755 scripts/entrypoint.sh /entrypoint.sh

COPY nsv_*.apk /tmp/
RUN apk add --no-cache --allow-untrusted /tmp/nsv_*.apk && rm /tmp/nsv_*.apk
