FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.8.0@sha256:add602d55daca18914838a78221f6bbe4284114b452c86a48f96d59aeb00f5c6 AS xx

FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine3.22@sha256:3587db7cc96576822c606d119729370dbf581931c5f43ac6d3fa03ab4ed85a10 AS builder

COPY --from=xx / /

RUN apk add --update --no-cache ca-certificates make git curl clang lld

ARG TARGETPLATFORM

RUN xx-apk --update --no-cache add musl-dev gcc

RUN xx-go --wrap

WORKDIR /usr/local/src/vault-operator

ARG GOPROXY

ENV CGO_ENABLED=0

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/vault-operator ./cmd/
RUN xx-verify /usr/local/bin/vault-operator


FROM alpine:3.23.0@sha256:51183f2cfa6320055da30872f211093f9ff1d3cf06f39a0bdb212314c5dc7375

RUN apk add --update --no-cache ca-certificates tzdata

COPY --from=builder /usr/local/bin/vault-operator /usr/local/bin/vault-operator

USER 65534

ENTRYPOINT ["vault-operator"]
