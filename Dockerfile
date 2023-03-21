# 构建
FROM golang:1.19 AS buildState
MAINTAINER baiyz0825<byz0825@outlook.com>
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"
WORKDIR /apps
COPY . /apps
RUN cd /apps && go build -o bot

# 打包
FROM alpine:latest
WORKDIR /apps
COPY --from=buildState /apps/bot /apps/
COPY --from=buildState /apps/config/config.yaml.example /apps/config/
# 设置时区为上海
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
# 设置编码
ENV LANG C.UTF-8
# 设置卷
VOLUME ["/apps/config"]
# 暴露端口
EXPOSE 50008 7890
ENTRYPOINT ["/apps/bot"]
