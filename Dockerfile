# Builder
FROM --platform=$BUILDPLATFORM whatwewant/builder-go:v1.25-1 AS builder

WORKDIR /build

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
  GOOS=$TARGETOS \
  GOARCH=$TARGETARCH \
  go build \
  -trimpath \
  -ldflags '-w -s -buildid=' \
  -v -o dns ./cmd/dns

# Server
FROM whatwewant/alpine:v3-1

LABEL MAINTAINER="Zero<tobewhatwewant@gmail.com>"

LABEL org.opencontainers.image.source="https://github.com/go-idp/dns"

COPY --from=builder /build/dns /bin

CMD /bin/dns server
