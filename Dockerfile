FROM golang:1.20 AS builder
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux go build -o omniplug


FROM kong:3.4.2

COPY --from=builder /src/omniplug /usr/local/bin/omniplug
COPY kong.conf /tmp
RUN cat /tmp/kong.conf >> /etc/kong/kong.conf 