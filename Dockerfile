FROM golang:1.19
MAINTAINER baiyz0825<byz0825@outlook.com>

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"

WORKDIR /apps
VOLUME ["/apps/config"]

COPY ./wxbot /apps/bin/wxbot
COPY ./config /apps/config

# 设置时区为上海
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# 设置编码
ENV LANG C.UTF-8
ENTRYPOINT ["/apps/bin/wxbot"]

