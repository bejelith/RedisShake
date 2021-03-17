FROM golang:1.14-buster
# Assumes gcc and bash are present
# RUN apk add git bash gcc

ARG goproxy=""
ENV GOPROXY=$goproxy
COPY ./ .
RUN /bin/bash -x ./build.sh


FROM golang:1.15-alpine

COPY --from=0 /go/bin/redis-shake.linux /usr/local/app/redis-shake
COPY --from=0 /go/wdconfig.conf /usr/local/app/redis-shake.conf
CMD /usr/local/app/redis-shake -type=sync -conf=/usr/local/app/redis-shake.conf
