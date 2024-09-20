FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /usr/src/app/

COPY ../ ./

RUN go mod download

RUN go build -o /go/bin/shop/ ./cmd/shop

# -----------------------------------------------------------------------------

FROM alpine:3.20.3

WORKDIR /usr/src/app/

COPY --from=builder /go/bin/shop/ ./

EXPOSE ${SHOP_SERVICE_PORT}
ENTRYPOINT [ "./shop" ]
