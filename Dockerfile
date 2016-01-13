FROM alpine
MAINTAINER Brett Fowle <brettfowle@gmail.com>

ENV BUILD_PATH /go/src/github.com/bfowle/avast
ENV BUILD_DEPS go git gcc libc-dev libgcc
ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go:$BUILD_PATH/vendor

RUN apk --update add $BUILD_DEPS

WORKDIR $BUILD_PATH
COPY . $BUILD_PATH

RUN go build -o /usr/bin/avast github.com/bfowle/avast/src/

#RUN apk --update add $BUILD_DEPS && \
#  cd $BUILD_PATH; go build -o /usr/bin/avast . && \
#  apk del $BUILD_DEPS && \
#  rm -rf /go /var/cache/apk/*

ENTRYPOINT ["avast"]
