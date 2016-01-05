FROM alpine
MAINTAINER Brett Fowle <brettfowle@gmail.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go:/go/src/github.com/bfowle/docker-hack/vendor

RUN apk update && \
  apk add ca-certificates && \
  rm -rf /var/cache/apk/*

COPY . /go/src/github.com/bfowle/docker-hack

RUN buildDeps='go \
  git \
  gcc \
  libc-dev \
  libgcc' \
  set -x && \
  apk update && \
  apk add $buildDeps && \
  cd /go/src/github.com/bfowle/docker-hack && \
  go build -o /usr/bin/docker-hack . && \
  apk del $buildDeps && \
  rm -rf /go /var/cache/apk/* && \
  echo "Build complete."

ENTRYPOINT ["docker-hack"]
