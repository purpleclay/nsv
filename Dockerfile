# syntax=docker/dockerfile:1
FROM alpine:3.20

RUN apk add --no-cache git git-lfs gnupg tini curl
RUN sh -c "$(curl https://raw.githubusercontent.com/purpleclay/gpg-import/main/scripts/install)" -- -v 0.3.2

ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
CMD ["--help"]

COPY --chmod=755 scripts/entrypoint.sh /entrypoint.sh

COPY nsv_*.apk /tmp/
RUN apk add --no-cache --allow-untrusted /tmp/nsv_*.apk && rm /tmp/nsv_*.apk
