FROM golang:1.17.5-alpine AS builder

WORKDIR /build/tsubaki
COPY . .
RUN go get
RUN make build

FROM alpine:latest

WORKDIR /opt/arisu/tsubaki
COPY --from=builder /build/tsubaki/build/tsubaki /opt/arisu/tsubaki/tsubaki

ENTRYPOINT ["sh", "/opt/arisu/tsubaki/tsubaki"]
