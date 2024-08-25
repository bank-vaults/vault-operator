FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.5.0@sha256:0c6a569797744e45955f39d4f7538ac344bfb7ebf0a54006a0a4297b153ccf0f AS xx

FROM --platform=$BUILDPLATFORM golang:1.22.3-alpine3.18@sha256:d1a601b64de09e2fa38c95e55838961811d5ca11062a8f4230a5c434b3ae2a34 AS builder

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


FROM alpine:3.20.2@sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5

RUN apk add --update --no-cache ca-certificates tzdata

COPY --from=builder /usr/local/bin/vault-operator /usr/local/bin/vault-operator

USER 65534

ENTRYPOINT ["vault-operator"]
