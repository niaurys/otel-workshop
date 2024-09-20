FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /usr/src/app/

COPY ../ ./

RUN go mod download

RUN go build -o /go/bin/factory/ ./cmd/factory

# -----------------------------------------------------------------------------

FROM alpine:3.20.3

WORKDIR /usr/src/app/

COPY --from=builder /go/bin/factory/ ./

ENTRYPOINT [ "./factory" ]
