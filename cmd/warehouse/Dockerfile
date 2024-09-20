FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /usr/src/app/

COPY ../ ./

RUN go mod download

RUN go build -o /go/bin/warehouse/ ./cmd/warehouse

# -----------------------------------------------------------------------------

FROM alpine:3.20.3

WORKDIR /usr/src/app/

COPY --from=builder /go/bin/warehouse/ ./

ENTRYPOINT [ "./warehouse" ]
