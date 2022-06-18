FROM golang:1.18.3-alpine AS builder

RUN apk update && apk add make git jq bash
WORKDIR /build/tsubaki
COPY . .
RUN go get
RUN make build

FROM alpine:latest

RUN apk update && apk add bash
WORKDIR /app/arisu/tsubaki
COPY --from=builder /build/tsubaki/bin/tsubaki                 /app/arisu/tsubaki/tsubaki
COPY --from=builder /build/tsubaki/docker/run.sh               /app/arisu/tsubaki/scripts/run.sh
COPY --from=builder /build/tsubaki/assets/banner.txt           /app/arisu/tsubaki/assets/banner.txt
COPY --from=builder /build/tsubaki/docker/lib/liblog.sh        /app/arisu/tsubaki/lib/liblog.sh
COPY --from=builder /build/tsubaki/docker/docker-entrypoint.sh /app/arisu/tsubaki/scripts/docker-entrypoint.sh

# Create a symbolic link to link it so you can do `docker run --rm arisuland/tsubaki:latest tsubaki version --json`
RUN ln -s /app/arisu/tsubaki/tsubaki /usr/bin/tsubaki

# not root.
USER 1001

# This is automatically set to production so you don't have to.
ENV GO_ENV=production

# This will be able to execute `tsubaki <command>`
ENTRYPOINT ["/app/arisu/tsubaki/scripts/docker-entrypoint.sh"]
CMD ["/app/arisu/tsubaki/scripts/run.sh"]
