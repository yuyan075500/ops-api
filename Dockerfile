FROM swr.cn-east-3.myhuaweicloud.com/lansefenbao/golang:1.23.1

# 指定工作目录
WORKDIR /app

# 复制源到工作目录
COPY . /app

# Go环境设置
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"

# 更改时区
RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo 'Asia/Shanghai' > /etc/timezone

# 处理项目依赖
RUN go mod tidy

# 项目编译
RUN go build

# 暴露端口
EXPOSE 8000

# 执行
ENTRYPOINT ["./ops-api"]
