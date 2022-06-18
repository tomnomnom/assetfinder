# Build
FROM golang:1.18.3-alpine AS build-env
RUN go install github.com/tomnomnom/assetfinder@latest

# Release
FROM alpine:3.16.0
RUN apk -U upgrade --no-cache \
    && apk add --no-cache bind-tools ca-certificates
COPY --from=build-env /go/bin/assetfinder /usr/local/bin/assetfinder

ENTRYPOINT ["assetfinder"]
