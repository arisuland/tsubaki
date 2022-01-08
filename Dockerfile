FROM golang:1.17.5-alpine AS builder

WORKDIR /build/tsubaki
COPY . .
RUN go get
RUN make build

FROM alpine:latest

WORKDIR /opt/arisu/tsubaki
COPY --from=builder /build/tsubaki/bin/tsubaki                 /opt/arisu/tsubaki/tsubaki
COPY --from=builder /build/tsubaki/docker/lib/liblog.sh        /opt/arisu/tsubaki/lib/liblog.sh
COPY --from=builder /build/tsubaki/docker/docker-entrypoint.sh /opt/arisu/tsubaki/docker-entrypoint.sh

ENTRYPOINT ["sh", "/opt/arisu/tsubaki/tsubaki"]
