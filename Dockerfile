FROM docker-dev-artifactory.workday.com/dpm/golang:1.14-alpine-gcc
# Assumes gcc and bash are present
# RUN apk add git bash gcc

COPY . .
ARG goproxy=""
ENV GOPROXY=$goproxy
ENV https_proxy=$goproxy
RUN git config http.proxy "$goproxy"
RUN cd src/vendor && govendor sync && cd ../..
RUN ./fastbuild.sh

FROM docker-dev-artifactory.workday.com/dpm/centos:7.7.1908

COPY --from=0 /go/bin/redis-shake.linux /usr/local/app/redis-shake
COPY --from=0 /go/wdconfig.conf /usr/local/app/redis-shake.conf
CMD /usr/local/app/redis-shake -type=sync -conf=/usr/local/app/redis-shake.conf
