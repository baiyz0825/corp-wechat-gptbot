# 构建
FROM golang:1.19 AS buildState
LABEL maintainer="baiyz0825<byz0825@outlook.com>"
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"
WORKDIR /apps
COPY . /apps
#RUN go install github.com/go-delve/delve/cmd/dlv@latest
# go dlv调试
#RUN cd /apps && go build  -gcflags="all=-N -l" -o bot
RUN cd /apps && go build  -o bot

# 打包
FROM ubuntu:latest AS env
WORKDIR /apps
COPY --from=buildState /apps/bot /apps/
COPY --from=buildState /apps/config/config.yaml.example /apps/config/
COPY --from=buildState /apps/assert /apps/assert/
# 拷贝go dlv调试
#COPY --from=buildState /go/bin/dlv /
# 设置时区为上海
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN sed -i 's/archive.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list
RUN apt-get update
RUN apt-get install -y wkhtmltopdf
# 设置cgo依赖
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
#RUN apk --no-cache add tzdata ca-certificates libc6-compat libgcc libstdc++
# 设置字体pdf转化乱码
RUN cp /apps/assert/simsun.ttc /usr/share/fonts
RUN mkdir /apps/db
RUN echo 'Asia/Shanghai' >/etc/timezone
# 设置编码
ENV LANG C.UTF-8
# 暴露端口
EXPOSE 50008 40000
ENTRYPOINT ["/apps/bot"]
# go dlv调试
#CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/apps/bot"]

