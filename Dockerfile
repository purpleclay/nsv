FROM cgr.dev/chainguard/wolfi-base

RUN apk add --no-cache git git-lfs tini

ENTRYPOINT ["/sbin/tini", "--", "/usr/bin/nsv"]

COPY nsv_*.apk /tmp/
RUN apk add --no-cache --allow-untrusted /tmp/nsv_*.apk && rm /tmp/nsv_*.apk
