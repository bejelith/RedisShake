FROM golang:1.14-alpine
RUN apk add git bash gcc
COPY . .

RUN ./fastbuild.sh

FROM alpine:3.12

COPY --from=0 /go/bin/redis-shake.linux /usr/local/app/redis-shake
COPY --from=0 /go/wdconfig.conf /usr/local/app/redis-shake.conf
CMD /usr/local/app/redis-shake -type=sync -conf=/usr/local/app/redis-shake.conf
