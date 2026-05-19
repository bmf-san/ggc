# syntax=docker/dockerfile:1
# This Dockerfile is consumed by GoReleaser; `ggc` is injected by the
# release pipeline, so building this image outside of GoReleaser will fail
# unless you place the pre-built binary next to this file.
FROM alpine:3.22@sha256:310c62b5e7ca5b08167e4384c68db0fd2905dd9c7493756d356e893909057601

# ggc shells out to `git`, so the runtime image must ship one. `ca-certificates`
# keeps HTTPS remote operations working. No other runtime dependencies.
RUN apk add --no-cache git ca-certificates \
    && addgroup -S ggc \
    && adduser -S -G ggc ggc

COPY ggc /usr/local/bin/ggc

USER ggc
WORKDIR /work

ENTRYPOINT ["/usr/local/bin/ggc"]
