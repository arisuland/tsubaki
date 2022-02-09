# Why is this still Alpine in release, and not using `scratch`?
# Because, Tsubaki's Docker scripts uses bash as the runtime, so that's
# pretty much why, so the image is usually around ~3.1-6.2mb, which is still pretty small!

FROM alpine:3.15

RUN apk update && apk add --no-cache bash musl-dev libc-dev gcompat
WORKDIR /app/arisu/tsubaki

COPY docker/docker-entrypoint.sh /app/arisu/tsubaki/scripts/docker-entrypoint.sh
COPY docker/lib/liblog.sh        /app/arisu/tsubaki/lib/liblog.sh
COPY docker/run.sh               /app/arisu/tsubaki/scripts/runner.sh
COPY tsubaki                     /app/arisu/tsubaki/tsubaki

# Create a symbolic link so you can use the `tsubaki` command globally,
# so you don't need to do `docker run --rm arisuland/tsubaki:latest /app/arisu/tsubaki/tsubaki <command>`,
# you can just do `docker run --rm arisuland/tsubaki:latest tsubaki <command>`!
RUN ln -s /app/arisu/tsubaki/tsubaki /usr/bin/tsubaki

# not root.
USER 1001

# Since we are required to have this, this is set to `production` automatically
# for you.
ENV GO_ENV=production

# This will be able to execute any command using `docker exec`
ENTRYPOINT ["/app/arisu/tsubaki/scripts/docker-entrypoint.sh"]

# This will run the server when using `docker run`
CMD ["/app/arisu/tsubaki/scripts/runner.sh"]
