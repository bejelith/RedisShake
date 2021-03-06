FROM docker-dev-artifactory.workday.com/dpm/golang:1.14-alpine-gcc
# Assumes gcc and bash are present
# RUN apk add git bash gcc

COPY . .
ARG goproxy=""
ENV GOPROXY=$goproxy

RUN ./build.sh

FROM docker-dev-artifactory.workday.com/dpm/golang:1.14-alpine

COPY --from=0 /go/bin/redis-shake.linux /usr/local/app/redis-shake
COPY --from=0 /go/wdconfig.conf /usr/local/app/redis-shake.conf
CMD /usr/local/app/redis-shake -type=sync -conf=/usr/local/app/redis-shake.conf
