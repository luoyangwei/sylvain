ARG REGISTRY=docker.io
FROM ${REGISTRY}/golang:1.22.1-alpine3.19 AS builder
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# ARG APP_ENV
# ENV APP_ENV=${APP_ENV:?err}

RUN export GOPROXY=https://goproxy.cn,direct && \
    go mod tidy && \
    go build -ldflags="-s -w" -o app

FROM ${REGISTRY}/alpine:3.16
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

COPY --from=builder . .

ENTRYPOINT [ "./app"]